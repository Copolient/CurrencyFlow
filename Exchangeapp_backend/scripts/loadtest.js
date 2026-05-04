import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// ============================================================
// 自定义指标
// ============================================================
const errorRate = new Rate('errors');
const latencyP99 = new Trend('latency_p99');
const requestCount = new Counter('total_requests');

// ============================================================
// 场景配置：阶梯式压力测试
// ============================================================
export const options = {
  scenarios: {
    // 场景 1：阶梯式爬升 → 维持 → 冲击
    ramp_up: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 1000 },    // 2 分钟爬升到 1000 VU
        { duration: '3m', target: 3000 },    // 3 分钟爬升到 3000 VU
        { duration: '5m', target: 5000 },    // 5 分钟爬升到 5000 VU
        { duration: '10m', target: 5000 },   // 维持 5000 VU 10 分钟
        { duration: '30s', target: 10000 },  // 30 秒内冲击到 10000 VU
        { duration: '5m', target: 10000 },   // 维持 10000 VU 5 分钟
        { duration: '2m', target: 0 },       // 2 分钟降压
      ],
      gracefulRampDown: '30s',
    },

    // 场景 2：恒定速率的健康检查
    health_check: {
      executor: 'constant-arrival-rate',
      rate: 100,
      timeUnit: '1s',
      duration: '27m30s',
      preAllocatedVUs: 200,
      exec: 'healthCheck',
    },
  },

  // ============================================================
  // 阈值断言：错误率 > 1% 自动停止
  // ============================================================
  thresholds: {
    'errors': [
      { threshold: 'rate<0.01', abortOnFail: true },   // 错误率 < 1%，否则中止
    ],
    'http_req_duration': [
      { threshold: 'p(95)<500', abortOnFail: false },   // P95 < 500ms（告警但不中止）
      { threshold: 'p(99)<2000', abortOnFail: false },  // P99 < 2s
    ],
    'http_req_failed': [
      { threshold: 'rate<0.01', abortOnFail: true },    // 失败率 < 1%
    ],
    'latency_p99': [
      { threshold: 'p(99)<2000' },
    ],
  },
};

// ============================================================
// 测试数据
// ============================================================
const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

const TEST_USER = {
  username: `loadtest_${Math.random().toString(36).slice(2, 8)}`,
  password: 'LoadTest@2024!',
};

const ARTICLES = [
  { title: 'Go Concurrency Patterns', content: 'Goroutines and channels...', preview: 'go-concurrency' },
  { title: 'Kubernetes Autoscaling', content: 'HPA and VPA deep dive...', preview: 'k8s-autoscaling' },
  { title: 'Redis Cluster Best Practices', content: 'Sharding and replication...', preview: 'redis-cluster' },
  { title: 'Observability in Production', content: 'Tracing, metrics, logging...', preview: 'observability' },
];

const CURRENCY_PAIRS = [
  { fromCurrency: 'USD', toCurrency: 'CNY', rate: 7.24 },
  { fromCurrency: 'EUR', toCurrency: 'USD', rate: 1.08 },
  { fromCurrency: 'GBP', toCurrency: 'JPY', rate: 190.5 },
  { fromCurrency: 'USD', toCurrency: 'EUR', rate: 0.92 },
];

// ============================================================
// Setup：注册测试用户并获取 Token
// ============================================================
export function setup() {
  const registerRes = http.post(`${BASE_URL}/api/auth/register`, TEST_USER);
  let token = '';

  if (registerRes.status === 200) {
    token = registerRes.json('token');
  } else {
    // 用户可能已存在，尝试登录
    const loginRes = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
      username: TEST_USER.username,
      password: TEST_USER.password,
    }), { headers: { 'Content-Type': 'application/json' } });
    if (loginRes.status === 200) {
      token = loginRes.json('token');
    }
  }

  return { token };
}

// ============================================================
// 主场景：混合业务流量
// ============================================================
export default function (data) {
  const token = data.token;
  const authHeaders = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': token,
    },
  };

  group('Public API', function () {
    // 1. 获取汇率（读操作，占比最高）
    group('Get Exchange Rates', function () {
      const res = http.get(`${BASE_URL}/api/exchangeRates`);
      const success = check(res, {
        'exchange rates status 200': (r) => r.status === 200,
        'exchange rates is array': (r) => Array.isArray(r.json()),
      });
      errorRate.add(!success);
      requestCount.add(1);
      latencyP99.add(res.timings.duration);
    });

    // 2. 获取文章列表（带缓存的读操作）
    group('Get Articles', function () {
      const res = http.get(`${BASE_URL}/api/articles`, authHeaders);
      const success = check(res, {
        'articles status 200': (r) => r.status === 200,
      });
      errorRate.add(!success);
      requestCount.add(1);
      latencyP99.add(res.timings.duration);
    });
  });

  group('Write Operations', function () {
    // 3. 创建文章（写操作）
    const article = randomItem(ARTICLES);
    group('Create Article', function () {
      const res = http.post(`${BASE_URL}/api/articles`, JSON.stringify(article), authHeaders);
      const success = check(res, {
        'create article status 201': (r) => r.status === 201,
      });
      errorRate.add(!success);
      requestCount.add(1);
    });

    // 4. 创建汇率
    const pair = randomItem(CURRENCY_PAIRS);
    group('Create Exchange Rate', function () {
      const res = http.post(`${BASE_URL}/api/exchangeRates`, JSON.stringify(pair), authHeaders);
      const success = check(res, {
        'create rate status 201': (r) => r.status === 201,
      });
      errorRate.add(!success);
      requestCount.add(1);
    });
  });

  group('Like Operations', function () {
    // 5. 点赞（Redis INCR）
    const articleId = '1';
    group('Like Article', function () {
      const res = http.post(`${BASE_URL}/api/articles/${articleId}/like`, null, authHeaders);
      const success = check(res, {
        'like status 200': (r) => r.status === 200,
      });
      errorRate.add(!success);
      requestCount.add(1);
    });

    // 6. 获取点赞数
    group('Get Likes', function () {
      const res = http.get(`${BASE_URL}/api/articles/${articleId}/like`, authHeaders);
      const success = check(res, {
        'likes status 200': (r) => r.status === 200,
      });
      errorRate.add(!success);
      requestCount.add(1);
    });
  });

  sleep(Math.random() * 0.5 + 0.1); // 100-600ms 随机思考时间
}

// ============================================================
// 健康检查场景
// ============================================================
export function healthCheck() {
  const res = http.get(`${BASE_URL}/healthz`);
  check(res, {
    'health check ok': (r) => r.status === 200,
  });
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	"exchangeapp/pkg/cache"
)

type AIAnalystService struct {
	rateHistoryRepo repository.RateHistoryRepository
	cache           cache.Cache
}

type AnalysisResult struct {
	Analysis     string      `json:"analysis"`
	Trend        string      `json:"trend"`
	KeyLevels    KeyLevels   `json:"keyLevels"`
	RiskWarning  string      `json:"riskWarning"`
}

type KeyLevels struct {
	Support  float64 `json:"support"`
	Resistance float64 `json:"resistance"`
}

func NewAIAnalystService(rateHistoryRepo repository.RateHistoryRepository, cache cache.Cache) *AIAnalystService {
	return &AIAnalystService{
		rateHistoryRepo: rateHistoryRepo,
		cache:           cache,
	}
}

func (s *AIAnalystService) Analyze(ctx context.Context, from, to, question string) (*AnalysisResult, error) {
	cacheKey := fmt.Sprintf("ai_analysis:%s:%s", from, to)

	// Check cache
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			var result AnalysisResult
			if json.Unmarshal([]byte(cached), &result) == nil {
				return &result, nil
			}
		}
	}

	// Get historical data
	end := time.Now()
	start := end.Add(-30 * 24 * time.Hour)
	histories, err := s.rateHistoryRepo.FindByPairAndTimeRange(from, to, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate history: %w", err)
	}

	if len(histories) == 0 {
		return &AnalysisResult{
			Analysis:    fmt.Sprintf("暂无 %s/%s 的历史数据，无法进行分析。", from, to),
			Trend:       "neutral",
			KeyLevels:   KeyLevels{},
			RiskWarning: "以上分析仅供参考，不构成投资建议",
		}, nil
	}

	// Calculate basic statistics
	result := s.calculateAnalysis(from, to, question, histories)

	// Cache result
	if s.cache != nil {
		data, _ := json.Marshal(result)
		_ = s.cache.Set(ctx, cacheKey, string(data), time.Hour)
	}

	return result, nil
}

func (s *AIAnalystService) calculateAnalysis(from, to, question string, histories []model.ExchangeRateHistory) *AnalysisResult {
	if len(histories) == 0 {
		return &AnalysisResult{
			Analysis:    "数据不足",
			Trend:       "neutral",
			KeyLevels:   KeyLevels{},
			RiskWarning: "以上分析仅供参考，不构成投资建议",
		}
	}

	// Calculate basic stats
	first := histories[0].Rate
	last := histories[len(histories)-1].Rate
	min := first
	max := first
	sum := 0.0

	for _, h := range histories {
		if h.Rate < min {
			min = h.Rate
		}
		if h.Rate > max {
			max = h.Rate
		}
		sum += h.Rate
	}

	avg := sum / float64(len(histories))
	change := (last - first) / first * 100

	// Determine trend
	trend := "neutral"
	if change > 1.0 {
		trend = "bullish"
	} else if change < -1.0 {
		trend = "bearish"
	}

	// Build analysis
	analysis := fmt.Sprintf("## %s/%s 近30天走势分析\n\n", from, to)
	analysis += fmt.Sprintf("**当前汇率:** %.4f\n\n", last)
	analysis += fmt.Sprintf("**30天涨跌幅:** %.2f%%\n\n", change)
	analysis += fmt.Sprintf("**30天最高:** %.4f\n", max)
	analysis += fmt.Sprintf("**30天最低:** %.4f\n", min)
	analysis += fmt.Sprintf("**30天均价:** %.4f\n\n", avg)

	switch trend {
	case "bullish":
		analysis += fmt.Sprintf("**趋势判断:** 看涨 (↑)\n\n%s兑%s近期呈现上涨趋势，累计上涨%.2f%%。", from, to, change)
	case "bearish":
		analysis += fmt.Sprintf("**趋势判断:** 看跌 (↓)\n\n%s兑%s近期呈现下跌趋势，累计下跌%.2f%%。", from, to, -change)
	default:
		analysis += fmt.Sprintf("**趋势判断:** 震荡 (→)\n\n%s兑%s近期在区间内震荡，波动幅度%.2f%%。", from, to, change)
	}

	if question != "" {
		analysis += fmt.Sprintf("\n\n**关于您的问题:** \"%s\"\n\n", question)
		analysis += "基于近期走势来看，"
		if trend == "bullish" {
			analysis += "上涨动能仍在，但需关注回调风险。建议设置好止损位。"
		} else if trend == "bearish" {
			analysis += "下跌趋势明显，建议观望或轻仓操作。关注下方支撑位。"
		} else {
			analysis += "市场处于震荡区间，建议在支撑位附近买入，阻力位附近卖出。"
		}
	}

	return &AnalysisResult{
		Analysis: analysis,
		Trend:    trend,
		KeyLevels: KeyLevels{
			Support:    min + (avg-min)*0.3,
			Resistance: max - (max-avg)*0.3,
		},
		RiskWarning: "以上分析仅供参考，不构成投资建议。汇率受多种因素影响，过去表现不代表未来走势。",
	}
}

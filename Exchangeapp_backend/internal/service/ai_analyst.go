package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"exchangeapp/internal/model"
	"exchangeapp/internal/repository"
	"exchangeapp/pkg/cache"
	"exchangeapp/pkg/config"
)

type AIAnalystService struct {
	rateHistoryRepo repository.RateHistoryRepository
	cache           cache.Cache
	llmCfg          config.LLMConfig
	httpClient      *http.Client
}

type AnalysisResult struct {
	Analysis    string    `json:"analysis"`
	Trend       string    `json:"trend"`
	KeyLevels   KeyLevels `json:"keyLevels"`
	RiskWarning string    `json:"riskWarning"`
}

type KeyLevels struct {
	Support    float64 `json:"support"`
	Resistance float64 `json:"resistance"`
}

// Anthropic Messages API request/response structs
type anthropicRequest struct {
	Model     string              `json:"model"`
	MaxTokens int                 `json:"max_tokens"`
	Messages  []anthropicMessage  `json:"messages"`
	System    string              `json:"system,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type llmAnalysisResponse struct {
	Analysis      string   `json:"analysis"`
	Trend         string   `json:"trend"`
	Support       float64  `json:"support"`
	Resistance    float64  `json:"resistance"`
	RiskWarning   string   `json:"risk_warning"`
	// Flexible fields for different LLM response formats
	Summary       string   `json:"analysis_summary"`
	Direction     string   `json:"change_30d_direction"`
	Outlook       string   `json:"future_outlook"`
	Reasons       []string `json:"possible_reasons"`
}

func NewAIAnalystService(rateHistoryRepo repository.RateHistoryRepository, cache cache.Cache, llmCfg config.LLMConfig) *AIAnalystService {
	return &AIAnalystService{
		rateHistoryRepo: rateHistoryRepo,
		cache:           cache,
		llmCfg:          llmCfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func sanitizeCacheKey(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *AIAnalystService) Analyze(ctx context.Context, from, to, question string) (*AnalysisResult, error) {
	// Use sanitized cache key to prevent injection
	cacheKey := "ai_analysis:" + sanitizeCacheKey(from, to, question)

	// Check cache (only for no-question requests to avoid caching generic answers)
	if question == "" && s.cache != nil {
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

	// Try LLM analysis first
	var result *AnalysisResult
	if s.llmCfg.APIKey != "" {
		result, err = s.analyzeWithLLM(ctx, from, to, question, histories)
		if err != nil {
			log.Printf("LLM analysis failed, falling back to local: %v", err)
		}
	}

	// Fallback to local calculation
	if result == nil {
		result = s.calculateAnalysis(from, to, question, histories)
	}

	// Cache result
	if s.cache != nil && question == "" {
		data, _ := json.Marshal(result)
		_ = s.cache.Set(ctx, cacheKey, string(data), time.Hour)
	}

	return result, nil
}

func (s *AIAnalystService) analyzeWithLLM(ctx context.Context, from, to, question string, histories []model.ExchangeRateHistory) (*AnalysisResult, error) {
	// Build statistics summary for the prompt
	stats := buildStatsSummary(from, to, histories)

	systemPrompt := `你是一位专业的外汇分析师。请根据提供的汇率历史数据进行分析。
请以JSON格式返回分析结果，格式如下：
{
  "analysis": "详细的分析文本（使用markdown格式）",
  "trend": "bullish/bearish/neutral",
  "support": 支撑价位(float64),
  "resistance": 阻力价位(float64),
  "risk_warning": "风险提示"
}
只返回JSON，不要有其他内容。`

	userPrompt := stats
	if question != "" {
		userPrompt += "\n\n用户问题：" + question
	}

	// Call Anthropic Messages API
	reqBody := anthropicRequest{
		Model:     s.llmCfg.Model,
		MaxTokens: s.llmCfg.MaxTokens,
		System:    systemPrompt,
		Messages: []anthropicMessage{
			{Role: "user", Content: userPrompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := strings.TrimRight(s.llmCfg.BaseURL, "/") + "/v1/messages"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.llmCfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Limit response body to 10MB to prevent memory exhaustion
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	// Find the text content block (skip thinking blocks)
	var text string
	for _, block := range apiResp.Content {
		if block.Type == "text" {
			text = block.Text
			break
		}
	}
	if text == "" {
		return nil, fmt.Errorf("no text content in LLM response")
	}
	// Extract JSON from possible markdown code blocks
	text = extractJSON(text)

	var llmResp llmAnalysisResponse
	if err := json.Unmarshal([]byte(text), &llmResp); err != nil {
		return nil, fmt.Errorf("parse LLM response: %w", err)
	}

	// Normalize trend from various formats
	trend := llmResp.Trend
	if trend == "" {
		switch llmResp.Direction {
		case "decline", "down", "bearish":
			trend = "bearish"
		case "rise", "up", "bullish":
			trend = "bullish"
		default:
			trend = "neutral"
		}
	}
	switch trend {
	case "bullish", "bearish", "neutral":
		// ok
	default:
		trend = "neutral"
	}

	// Build analysis from whatever fields the LLM returned
	analysis := llmResp.Analysis
	if analysis == "" && llmResp.Summary != "" {
		analysis = llmResp.Summary
		if llmResp.Outlook != "" {
			analysis += "\n\n**未来展望:** " + llmResp.Outlook
		}
		if len(llmResp.Reasons) > 0 {
			analysis += "\n\n**可能原因:**\n"
			for _, r := range llmResp.Reasons {
				analysis += "- " + r + "\n"
			}
		}
	}

	riskWarning := llmResp.RiskWarning
	if riskWarning == "" {
		riskWarning = "以上分析仅供参考，不构成投资建议。汇率受多种因素影响，过去表现不代表未来走势。"
	}

	return &AnalysisResult{
		Analysis: analysis,
		Trend:    trend,
		KeyLevels: KeyLevels{
			Support:    llmResp.Support,
			Resistance: llmResp.Resistance,
		},
		RiskWarning: riskWarning,
	}, nil
}

func buildStatsSummary(from, to string, histories []model.ExchangeRateHistory) string {
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

	// Recent 7 days data
	recentStart := len(histories) - 7*24
	if recentStart < 0 {
		recentStart = 0
	}
	recentFirst := histories[recentStart].Rate
	recentChange := (last - recentFirst) / recentFirst * 100

	return fmt.Sprintf(`请分析以下 %s/%s 汇率数据（近30天）：

- 当前汇率: %.4f
- 30天前汇率: %.4f
- 30天涨跌幅: %.2f%%
- 30天最高: %.4f
- 30天最低: %.4f
- 30天均价: %.4f
- 近7天涨跌幅: %.2f%%
- 数据点数量: %d
- 时间范围: %s 至 %s`,
		from, to,
		last, first, change,
		max, min, avg,
		recentChange,
		len(histories),
		histories[0].Timestamp.Format("2006-01-02"),
		histories[len(histories)-1].Timestamp.Format("2006-01-02"),
	)
}

func extractJSON(text string) string {
	// Try to extract JSON from markdown code blocks
	if idx := strings.Index(text, "```json"); idx != -1 {
		text = text[idx+7:]
		if endIdx := strings.Index(text, "```"); endIdx != -1 {
			return strings.TrimSpace(text[:endIdx])
		}
	}
	if idx := strings.Index(text, "```"); idx != -1 {
		text = text[idx+3:]
		if endIdx := strings.Index(text, "```"); endIdx != -1 {
			return strings.TrimSpace(text[:endIdx])
		}
	}
	// Try to find JSON object directly
	if idx := strings.Index(text, "{"); idx != -1 {
		if endIdx := strings.LastIndex(text, "}"); endIdx > idx {
			return text[idx : endIdx+1]
		}
	}
	return strings.TrimSpace(text)
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

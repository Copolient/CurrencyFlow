package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"exchangeapp/internal/model"
	"exchangeapp/pkg/config"
	"exchangeapp/pkg/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	db := database.NewMySQL(cfg.Database)

	log.Println("Seeding database...")

	seedUsers(db)
	seedExchangeRates(db)
	seedExchangeRateHistories(db)
	seedArticles(db)
	seedPosts(db)
	seedFavorites(db)
	seedFollows(db)
	seedRateAlerts(db)
	seedNotifications(db)

	log.Println("Seed completed!")
}

func seedUsers(db *gorm.DB) {
	users := []struct {
		Username string
		Bio      string
		Avatar   string
	}{
		{"alice", "外汇交易爱好者，专注美元兑人民币走势分析", "https://api.dicebear.com/7.x/avataaars/svg?seed=alice"},
		{"bob", "量化交易分析师，10年金融从业经验", "https://api.dicebear.com/7.x/avataaars/svg?seed=bob"},
		{"charlie", "财经自媒体博主，每日汇率解读", "https://api.dicebear.com/7.x/avataaars/svg?seed=charlie"},
		{"diana", "留学党，关注英镑和日元汇率", "https://api.dicebear.com/7.x/avataaars/svg?seed=diana"},
		{"eve", "跨境电商从业者，实时关注汇率变动", "https://api.dicebear.com/7.x/avataaars/svg?seed=eve"},
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 10)

	for _, u := range users {
		user := model.User{
			Username: u.Username,
			Password: string(hash),
			Bio:      u.Bio,
			Avatar:   u.Avatar,
		}
		db.Where("username = ?", u.Username).FirstOrCreate(&user)
	}
	log.Println("  Users seeded (5)")
}

func seedExchangeRates(db *gorm.DB) {
	pairs := []struct{ From, To, Rate string }{
		{"USD", "CNY", "7.2450"},
		{"USD", "EUR", "0.9210"},
		{"USD", "GBP", "0.7920"},
		{"USD", "JPY", "149.50"},
		{"EUR", "CNY", "7.8650"},
		{"EUR", "GBP", "0.8600"},
		{"EUR", "JPY", "162.30"},
		{"GBP", "CNY", "9.1420"},
		{"GBP", "JPY", "188.75"},
		{"CNY", "JPY", "20.65"},
	}

	for _, p := range pairs {
		rate := parseRate(p.Rate)
		er := model.ExchangeRate{
			FromCurrency: p.From,
			ToCurrency:   p.To,
			Rate:         rate,
			Date:         time.Now(),
		}
		db.Where("from_currency = ? AND to_currency = ?", p.From, p.To).FirstOrCreate(&er)
	}
	log.Println("  Exchange rates seeded (10 pairs)")
}

func seedExchangeRateHistories(db *gorm.DB) {
	baseRates := map[string]float64{
		"USD_CNY": 7.2450, "USD_EUR": 0.9210, "USD_GBP": 0.7920, "USD_JPY": 149.50,
		"EUR_CNY": 7.8650, "EUR_GBP": 0.8600, "EUR_JPY": 162.30,
		"GBP_CNY": 9.1420, "GBP_JPY": 188.75, "CNY_JPY": 20.65,
	}

	now := time.Now()
	count := 0

	for pair, baseRate := range baseRates {
		from := pair[:3]
		to := pair[4:]

		// Generate 30 days of data, 1 record per hour = 720 records per pair
		for d := 0; d < 30; d++ {
			for h := 0; h < 24; h++ {
				ts := now.Add(-time.Duration(30-d) * 24 * time.Hour).Add(-time.Duration(24-h) * time.Hour)

				// Simulate realistic price movement with drift and volatility
				daysFromStart := float64(d)
				drift := (rand.Float64() - 0.48) * 0.001 * daysFromStart // slight upward drift
				volatility := (rand.Float64() - 0.5) * baseRate * 0.008  // ~0.8% hourly volatility
				rate := baseRate + drift + volatility

				// Ensure rate stays positive and reasonable
				if rate <= 0 {
					rate = baseRate * 0.95
				}

				history := model.ExchangeRateHistory{
					FromCurrency: from,
					ToCurrency:   to,
					Rate:         math.Round(rate*10000) / 10000,
					Timestamp:    ts,
				}
				db.Create(&history)
				count++
			}
		}
	}
	log.Printf("  Exchange rate histories seeded (%d records)\n", count)
}

func seedArticles(db *gorm.DB) {
	articles := []struct {
		Title, Content, Preview string
	}{
		{
			"美联储议息会议纪要释放鸽派信号，美元指数承压下行",
			"美联储最新公布的议息会议纪要显示，多数委员认为当前利率水平已经足够限制性，未来可能考虑降息。这一鸽派信号导致美元指数在纽约交易时段下跌0.3%，美元兑人民币汇率从7.25回落至7.24附近。分析人士指出，如果美联储在下半年开启降息周期，美元可能会进一步走弱。投资者需密切关注即将公布的非农就业数据，这将成为美联储决策的重要参考依据。",
			"美联储最新议息会议纪要释放鸽派信号，美元指数承压下行0.3%",
		},
		{
			"英国央行维持利率不变，英镑兑美元短线拉升50点",
			"英国央行在最新的货币政策会议上决定维持基准利率在5.25%不变，这一决定符合市场预期。然而，央行行长在新闻发布会上暗示，如果通胀持续回落，未来可能考虑降息。英镑兑美元在决议公布后短线拉升50个基点，最高触及0.7950。市场分析师认为，英国经济数据的好转为英镑提供了支撑，但脱欧后的结构性问题仍然限制了英镑的上行空间。建议关注下周公布的英国GDP数据。",
			"英国央行维持利率不变，英镑短线拉升50点，关注后续经济数据",
		},
		{
			"日本央行结束负利率时代，日元迎来历史性转折",
			"日本央行在最新的货币政策会议上做出了历史性的决定——结束长达8年的负利率政策，将基准利率上调至0%-0.1%区间。这一决定标志着日本货币政策的重大转向。美元兑日元在消息公布后大幅波动，从150快速下跌至148，随后又反弹至149.50附近。市场对于日本央行未来的加息路径存在分歧，部分分析师认为年内可能还有一次加息，而另一部分则认为央行会保持谨慎。对于持有日元资产的投资者来说，这是一个重要的里程碑。",
			"日本央行结束8年负利率时代，日元迎来历史性转折，美元兑日元大幅波动",
		},
		{
			"人民币国际化再提速，跨境人民币结算量创新高",
			"中国人民银行最新数据显示，2024年跨境人民币结算量突破50万亿元，同比增长35%，创历史新高。人民币在全球支付货币中的份额已升至第四位。分析认为，人民币国际化的加速推进，一方面得益于中国与一带一路沿线国家贸易的增长，另一方面也反映了全球去美元化的趋势。对于外汇市场而言，人民币国际化的深化有助于降低美元兑人民币汇率的波动性，增强人民币的定价权。投资者可关注人民币在全球储备货币中的占比变化。",
			"跨境人民币结算量突破50万亿元创新高，人民币国际化再提速",
		},
		{
			"欧元区经济复苏乏力，欧央行面临政策两难",
			"欧元区最新公布的PMI数据低于预期，制造业PMI连续15个月处于收缩区间，服务业PMI也出现下滑。欧洲央行面临两难局面：一方面通胀仍高于2%的目标水平，另一方面经济增长疲软。欧元兑美元在数据公布后下跌至0.9180附近。分析师指出，如果欧元区经济数据持续恶化，欧央行可能不得不提前降息，这将对欧元构成进一步压力。建议关注下周公布的欧元区CPI数据和欧央行官员讲话。",
			"欧元区PMI数据不佳，欧央行面临通胀与增长的两难抉择，欧元承压",
		},
	}

	for _, a := range articles {
		article := model.Article{
			Title:   a.Title,
			Content: a.Content,
			Preview: a.Preview,
		}
		db.Where("title = ?", a.Title).FirstOrCreate(&article)
	}
	log.Println("  Articles seeded (5)")
}

func seedPosts(db *gorm.DB) {
	var users []model.User
	db.Find(&users)
	if len(users) == 0 {
		log.Println("  Posts skipped (no users)")
		return
	}

	posts := []struct {
		Content  string
		Currency string
		UserIdx  int
	}{
		{"今天美元兑人民币跌破7.24了，感觉短期内还会继续下行，大家怎么看？", "USD/CNY", 0},
		{"刚换了1000英镑交学费，汇率比上个月好多了，开心！", "GBP/CNY", 3},
		{"日元终于要加息了！持有日元资产的朋友要注意了，可能会有一波行情。", "USD/JPY", 1},
		{"做跨境电商的朋友们注意，欧元最近波动比较大，建议锁定汇率风险。", "EUR/CNY", 4},
		{"美联储降息预期升温，黄金和非美货币都在涨，美元指数跌破104。", "USD", 2},
		{"分享一个汇率分析技巧：关注10年期美债收益率，它和美元指数高度相关。", "USD", 1},
		{"人民币国际化进程加速，长期来看美元兑人民币可能会回到7以下。", "USD/CNY", 0},
		{"英镑最近走势不错，英国央行维持高利率对英镑是利好。", "GBP/CNY", 2},
		{"日本央行结束负利率，这是历史性的时刻！日元可能迎来长期升值。", "USD/JPY", 3},
		{"建议大家设置汇率预警，不要错过最佳换汇时机。我已经设了好几个。", "USD/CNY", 4},
	}

	for _, p := range posts {
		post := model.Post{
			UserID:   users[p.UserIdx].ID,
			Content:  p.Content,
			Currency: p.Currency,
			Likes:    rand.Intn(50),
		}
		db.Create(&post)
	}
	log.Println("  Posts seeded (10)")
}

func seedFavorites(db *gorm.DB) {
	var users []model.User
	db.Find(&users)
	if len(users) < 3 {
		return
	}

	favs := []struct {
		UserIdx int
		From, To string
	}{
		{0, "USD", "CNY"},
		{0, "EUR", "CNY"},
		{1, "USD", "JPY"},
		{1, "GBP", "USD"},
		{3, "GBP", "CNY"},
		{3, "USD", "JPY"},
		{4, "EUR", "CNY"},
		{4, "USD", "CNY"},
	}

	for _, f := range favs {
		fav := model.Favorite{
			UserID:       users[f.UserIdx].ID,
			FromCurrency: f.From,
			ToCurrency:   f.To,
		}
		db.Where("user_id = ? AND from_currency = ? AND to_currency = ?", fav.UserID, f.From, f.To).FirstOrCreate(&fav)
	}
	log.Println("  Favorites seeded (8)")
}

func seedFollows(db *gorm.DB) {
	var users []model.User
	db.Find(&users)
	if len(users) < 5 {
		return
	}

	// alice(0) follows bob(1), charlie(2)
	// bob(1) follows alice(0), charlie(2), diana(3)
	// diana(3) follows alice(0), eve(4)
	// eve(4) follows bob(1), charlie(2)
	follows := []struct{ From, To int }{
		{0, 1}, {0, 2},
		{1, 0}, {1, 2}, {1, 3},
		{3, 0}, {3, 4},
		{4, 1}, {4, 2},
	}

	for _, f := range follows {
		follow := model.Follow{
			FollowerID: users[f.From].ID,
			FolloweeID: users[f.To].ID,
		}
		db.Where("follower_id = ? AND followee_id = ?", follow.FollowerID, follow.FolloweeID).FirstOrCreate(&follow)
	}

	// Update follower/following counts
	for _, u := range users {
		var followersCount, followingCount int64
		db.Model(&model.Follow{}).Where("followee_id = ?", u.ID).Count(&followersCount)
		db.Model(&model.Follow{}).Where("follower_id = ?", u.ID).Count(&followingCount)
		db.Model(&model.User{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
			"followers_count": followersCount,
			"following_count": followingCount,
		})
	}
	log.Println("  Follows seeded (9)")
}

func seedRateAlerts(db *gorm.DB) {
	var users []model.User
	db.Find(&users)
	if len(users) < 2 {
		return
	}

	alerts := []struct {
		UserIdx                        int
		From, To, Direction, TargetRate string
	}{
		{0, "USD", "CNY", "below", "7.2000"},
		{0, "EUR", "CNY", "above", "8.0000"},
		{1, "USD", "JPY", "above", "150.0000"},
		{3, "GBP", "CNY", "below", "9.0000"},
		{4, "USD", "CNY", "above", "7.3000"},
	}

	for _, a := range alerts {
		alert := model.RateAlert{
			UserID:       users[a.UserIdx].ID,
			FromCurrency: a.From,
			ToCurrency:   a.To,
			TargetRate:   parseRate(a.TargetRate),
			Direction:    a.Direction,
		}
		db.Create(&alert)
	}
	log.Println("  Rate alerts seeded (5)")
}

func seedNotifications(db *gorm.DB) {
	var users []model.User
	db.Find(&users)
	if len(users) < 2 {
		return
	}

	notifs := []struct {
		UserIdx        int
		Type, Title, Content string
		Read           bool
	}{
		{0, "alert_triggered", "汇率预警触发", "USD/CNY 已跌破 7.2000，当前汇率 7.1985", false},
		{0, "system", "欢迎加入 ExchangeApp", "感谢注册！您可以设置汇率预警、收藏常用货币对。", true},
		{1, "alert_triggered", "汇率预警触发", "USD/JPY 已突破 150.0000，当前汇率 150.25", false},
		{3, "system", "新功能上线", "AI 智能分析师已上线，快来体验 AI 驱动的汇率分析！", false},
		{4, "alert_triggered", "汇率预警触发", "USD/CNY 已突破 7.3000，当前汇率 7.3025", false},
	}

	for _, n := range notifs {
		notif := model.Notification{
			UserID:  users[n.UserIdx].ID,
			Type:    n.Type,
			Title:   n.Title,
			Content: n.Content,
			Read:    n.Read,
		}
		db.Create(&notif)
	}
	log.Println("  Notifications seeded (5)")
}

func parseRate(s string) float64 {
	var r float64
	fmt.Sscanf(s, "%f", &r)
	return r
}

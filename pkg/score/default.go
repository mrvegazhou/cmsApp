package score

import (
	"math"
	"time"
)

type Comment struct {
	Likes     int
	Dislikes  int
	Replies   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 威尔逊得分计算函数
func Wilson(ups, downs int) float64 {
	n := ups + downs
	if n == 0 {
		return 0
	}
	n1 := float64(n)
	z := 1.96
	p := float64(ups / n)
	zzfn := z * z / (4 * n1)
	return (p + 2.0*zzfn - z*math.Sqrt((zzfn/n1+p*(1.0-p))/n1)) / (1 + 4*zzfn)
}

// 时间衰减因子
func TimeDecayFactor(createdAt time.Time, updatedAt time.Time) float64 {
	now := time.Now()
	duration := now.Sub(createdAt)
	const timeFactor = 0.0005 // 衰减率，控制时间影响的程度
	qAge := duration.Hours() / 24
	qUpdated := (updatedAt.Sub(createdAt)).Hours() / 24
	elapsed := qAge - qUpdated/2
	return math.Exp(-timeFactor * elapsed)
}

// 计算综合得分
func (comment Comment) CalculateScoreByComment() float64 {
	// 威尔逊得分
	wilsonScoreValue := Wilson(comment.Likes, comment.Dislikes)

	// 评论回复数加权
	var normalizedReplies float64
	if comment.Replies <= 1 {
		// 使用 math.Log1p 处理接近 0 的值
		normalizedReplies = math.Log1p(float64(comment.Replies))
	} else {
		normalizedReplies = math.Log(float64(comment.Replies))
	}

	// 时间衰减因子
	timeDecay := TimeDecayFactor(comment.CreatedAt, comment.UpdatedAt)

	score := ((wilsonScoreValue + normalizedReplies/2) * timeDecay) / 2
	return score
}

type Article struct {
	Likes     int
	Dislikes  int
	Comments  int
	Views     int
	Shares    int
	Reports   int // 处罚数量
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (article Article) CalculateScoreByArticle() float64 {
	// 时间衰减因子
	timeDecay := TimeDecayFactor(article.CreatedAt, article.UpdatedAt)
	var normalizedComments float64
	if article.Comments <= 1 {
		// 使用 math.Log1p 处理接近 0 的值
		normalizedComments = math.Log1p(float64(article.Comments)) / 2
	} else {
		normalizedComments = math.Log(float64(article.Comments)) / 2
	}
	var normalizedViews float64
	if article.Views <= 1 {
		// 使用 math.Log1p 处理接近 0 的值
		normalizedViews = math.Log1p(float64(article.Views+article.Shares)) / 4
	} else {
		normalizedViews = math.Log(float64(article.Views+article.Shares)) / 4
	}

	// 威尔逊得分
	wilsonScoreValue := Wilson(article.Likes, article.Dislikes)

	score := ((wilsonScoreValue + normalizedComments + normalizedViews) * timeDecay) / 2
	return score
}

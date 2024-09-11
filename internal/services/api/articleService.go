package api

import (
	"cmsApp/configs"
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/utils/arrayx"
	"cmsApp/pkg/utils/imagex"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

type apiArticleService struct {
	Dao               *dao.AppArticleDao
	ImgsTempDao       *dao.ImgsTempDao
	ArticleHistoryDao *dao.AppArticleHistoryDao
}

var (
	instanceApiArticleService *apiArticleService
	onceApiArticleService     sync.Once
)

func NewApiArticleService() *apiArticleService {
	onceApiArticleService.Do(func() {
		instanceApiArticleService = &apiArticleService{
			Dao:               dao.NewAppArticleDao(),
			ImgsTempDao:       dao.NewImgsTempDao(),
			ArticleHistoryDao: dao.NewAppArticleHistoryDao(),
		}
	})
	return instanceApiArticleService
}

func (ser *apiArticleService) GetArticleInfo(id uint64) (article models.AppArticle, err error) {
	condition := map[string]interface{}{
		"id": id,
	}
	articleInfo, err := ser.Dao.GetAppArticle(condition)
	if err == gorm.ErrRecordNotFound {
		return models.AppArticle{}, nil
	}
	return articleInfo, err
}

// 检测一天内上传图片的总数
func (ser *apiArticleService) CheckUploadLimitNum(userId uint64) (bool, error) {
	conds := make(map[string][]interface{})
	exp := []interface{}{"= ?", userId}
	conds["user_id"] = exp
	exp = []interface{}{"= ?", time.Now()}
	conds["create_time"] = exp
	total, err := ser.ImgsTempDao.GetImgsTempTotal(conds)
	if err != nil {
		return false, errors.New(constant.IMAGE_UPLOAD_ERR)
	} else {
		if total > cast.ToInt64(configs.App.Upload.LimitNum) {
			return true, errors.New(constant.UPLOAD_EXCEED_ERR)
		} else {
			return false, nil
		}
	}
}

func (ser *apiArticleService) UploadImage(req models.AppImgTempUploadReq, userId uint64) (articleId uint64, fullPath string, imgName string, fileName string, err error) {
	fullPath = ""
	fileName = req.File.Filename
	resourceId := cast.ToUint64(req.ResourceId)
	if resourceId <= 0 {
		// 生成空文章信息
		article := models.AppArticle{}
		article.State = 1
		article.AuthorId = userId
		article.CreateTime = time.Now()
		article.UpdateTime = time.Now()
		articleId, err = ser.Dao.CreateAppArticle(article)
		if err != nil {
			return articleId, fullPath, imgName, fileName, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		req.ResourceId = cast.ToString(articleId)
	}
	// 存储到临时图片表
	_, _, imgName, fullPath, err = NewApiImgsTempService().SaveImage(req, userId)
	if err != nil {
		return resourceId, fullPath, imgName, fileName, err
	}
	return resourceId, fullPath, imgName, fileName, nil
}

func (ser *apiArticleService) UploadCoverImage(req models.AppImgTempUploadReq, userId uint64) (articleId uint64, fullPath string, imgName string, fileName string, err error) {
	resourceId := cast.ToUint64(req.ResourceId)
	if resourceId <= 0 {
		// 生成空文章信息
		article := models.AppArticle{}
		article.State = 1
		article.AuthorId = userId
		article.CreateTime = time.Now()
		article.UpdateTime = time.Now()
		articleId, err = ser.Dao.CreateAppArticle(article)
		if err != nil {
			return articleId, fullPath, imgName, fileName, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		req.ResourceId = cast.ToString(articleId)
	}
	_, _, imgName, fullPath, err = NewApiImgsTempService().SaveImage(req, userId)
	if err != nil {
		return resourceId, fullPath, imgName, fileName, err
	}
	return resourceId, fullPath, imgName, fileName, nil
}

// 保存草稿
func (ser *apiArticleService) SaveArticleDraft(userId uint64, req models.ArticleDraft) (article models.AppArticle, err error) {
	copier.CopyWithOption(&article, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if len(req.Tags) > 0 {
		article.Tags = arrayx.JoinIntArr2Str(req.Tags)
	}
	if req.AuthorId <= 0 {
		article.AuthorId = userId
	}
	if req.CoverUrl != "" {
		coverUrl := strings.Split(req.CoverUrl, "/")
		article.CoverUrl = coverUrl[len(coverUrl)-1]
	}
	if req.Marks != "" {
		article.Marks = req.Marks
	}
	article.State = 1
	article.UpdateTime = time.Now()
	// 保存文章
	if req.ArticleId != 0 {
		if req.ArticleId <= 0 {
			return article, errors.New(constant.ARTICLE_UPDATE_ERR)
		}
		// 检查是否存在文章
		articleInfo, err := ser.Dao.GetAppArticle(map[string]interface{}{
			"id": req.ArticleId,
		})
		if err == gorm.ErrRecordNotFound {
			return article, errors.New(constant.ARTICLE_UPDATE_ERR)
		}
		// 更新
		if err == nil {
			ser.Dao.UpdateArticle(req.ArticleId, article)
			article.Id = articleInfo.Id
		} else {
			return article, errors.New(constant.ARTICLE_UPDATE_ERR)
		}
	} else {
		article.CreateTime = time.Now()
		id, err := ser.Dao.CreateAppArticle(article)
		if err != nil {
			return article, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		article.Id = id
	}
	// 保存草稿到历史 判断内容是否相同，相同则更新时间，不同则增添一条新记录
	lastHistory, err := ser.ArticleHistoryDao.GetLastArticleHistoryInfo()
	if strings.Compare(lastHistory.Content, req.Content) == 0 && strings.Compare(lastHistory.Title, req.Title) == 0 {
		ser.ArticleHistoryDao.UpdateColumns(map[string]interface{}{
			"id": lastHistory.Id,
		}, map[string]interface{}{
			"update_time": time.Now(),
		}, nil)
	} else {
		var draftHistory models.AppArticleHistory
		copier.CopyWithOption(&draftHistory, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		draftHistory.Id = 0
		draftHistory.CreateTime = time.Now()
		draftHistory.Tags = article.Tags
		draftHistory.ArticleId = article.Id
		draftHistory.AuthorId = userId
		_, err = ser.ArticleHistoryDao.CreateArticleHistory(draftHistory)
		if err != nil {
			return article, errors.New(constant.ARTICLE_DRAFT_HISTORY_ERR)
		}
	}

	go func() {
		// 删除多余草稿
		ser.GetDraftHistoryInfo(article.Id)
	}()
	return article, nil
}

// 发布文章
func (ser *apiArticleService) SaveArticle(userId uint64, req models.Article) (article models.AppArticle, err error) {
	copier.CopyWithOption(&article, req, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if len(req.Tags) > 0 {
		article.Tags = arrayx.JoinIntArr2Str(req.Tags)
	}
	if req.AuthorId <= 0 {
		article.AuthorId = userId
	}
	if req.CoverUrl != "" {
		coverUrl := strings.Split(req.CoverUrl, "/")
		article.CoverUrl = coverUrl[len(coverUrl)-1]
	}
	if req.Marks != "" {
		article.Marks = req.Marks
	}
	article.State = 2
	now := time.Now()
	article.UpdateTime = now
	article.CreateTime = now
	// 保存文章
	if req.Id != 0 {
		// 检查是否存在文章
		articleInfo, err := ser.Dao.GetAppArticle(map[string]interface{}{
			"id": req.Id,
		})
		if err == gorm.ErrRecordNotFound {
			return article, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		// 更新
		if err == nil {
			ser.Dao.UpdateArticle(req.Id, article)
			article.Id = articleInfo.Id
		} else {
			return article, errors.New(constant.ARTICLE_SAVE_ERR)
		}
	} else {
		id, err := ser.Dao.CreateAppArticle(article)
		if err != nil {
			return article, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		article.Id = id
	}
	// 检查图片，删除多余
	imgs := imagex.ExtractImageSrcs(article.Content)
	err = NewApiImgsTempService().move2ArticleImgs(imgs, article.Id)
	if err != nil {
		return article, err
	}
	return article, nil
}

func (ser *apiArticleService) GetArticleList(articleIds []uint64) ([]models.AppArticle, error) {
	conditions := map[string][]interface{}{}
	if len(articleIds) > 0 {
		conditions = map[string][]interface{}{
			"id": {"IN ?", articleIds},
		}
		return ser.Dao.GetArticleList(conditions)
	} else {
		return []models.AppArticle{}, nil
	}
}

// 获取文章编辑的历史记录
func (ser *apiArticleService) GetDraftHistoryList(articleId uint64) ([]models.AppArticleHistory, error) {
	conditions := map[string][]interface{}{}
	if articleId != 0 {
		conditions = map[string][]interface{}{
			"article_id": {"= ?", articleId},
		}
		return ser.ArticleHistoryDao.GetArticleHistoryList(conditions)
	} else {
		return []models.AppArticleHistory{}, errors.New(constant.ARTICLE_DARFT_PARAM_ERR)
	}
}

// 判断草稿本是否满50条，超出删除多余的
func (ser *apiArticleService) DelExcessArticles(articleId uint64) {
	conditions := map[string][]interface{}{
		"article_id": {"= ?", articleId},
	}
	list, err := ser.ArticleHistoryDao.GetArticleHistoryList(conditions)
	if err == nil {
		if len(list) >= 50 {
			var ids []uint64
			for i, item := range list {
				if i > 50 {
					ids = append(ids, item.Id)
				}
			}
			// 删除
			conditions = map[string][]interface{}{
				"id": {"in ?", ids},
			}
			ser.ArticleHistoryDao.DelArticleHistory(conditions)
		}
	}
}

// 获取单条草稿详情
func (ser *apiArticleService) GetDraftHistoryInfo(id uint64) (models.ArticleHistoryResp, error) {
	historyInfo := models.ArticleHistoryResp{}
	if id != 0 {
		info, err := ser.ArticleHistoryDao.GetArticleHistoryInfo(map[string]interface{}{"id": id})
		if err != nil {
			return historyInfo, errors.New(constant.ARTICLE_DRAFT_HISTORY_ERR)
		}
		ids := arrayx.String2Uint64(strings.Split(info.Tags, ","))
		tagList, err := NewApiTagService().GetTagListByIds(ids)
		tagArr := make([]models.AppTagInfo, 0, len(tagList))
		for _, item := range tagList {
			tagArr = append(tagArr, models.AppTagInfo{Id: item.Id, Name: item.Name})
		}
		copier.CopyWithOption(&historyInfo, info, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		historyInfo.Tags = tagArr
		return historyInfo, nil
	} else {
		return historyInfo, errors.New(constant.ARTICLE_DARFT_PARAM_ERR)
	}
}

// 更新文章的评论总数
func (ser *apiArticleService) UpdateArticleCommentCount(id uint64) error {
	if &id == nil {
		return errors.New(constant.ARTICLE_UPDATE_COMMENT_COUNT_ERR)
	}
	return ser.Dao.UpdateColumnCount(id, "comment_count", "+", nil)
}

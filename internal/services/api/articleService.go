package api

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"errors"
	"gorm.io/gorm"
	"net/url"
	"sync"
	"time"
)

type apiArticleService struct {
	Dao *dao.AppArticleDao
}

var (
	instanceApiArticleService *apiArticleService
	onceApiArticleService     sync.Once
)

func NewApiArticleService() *apiArticleService {
	onceApiArticleService.Do(func() {
		instanceApiArticleService = &apiArticleService{
			Dao: dao.NewAppArticleDao(),
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

func (ser *apiArticleService) UploadImage(req models.AppArticleUploadImage, userId uint64) (fullPath string, imgName string, fileName string, err error) {
	fullPath = ""
	fileName = req.File.Filename
	if req.ArticleId <= 0 {
		// 生成空文章信息
		article := models.AppArticle{}
		article.State = 1
		article.AuthorId = userId
		article.CreateTime = time.Now()
		article.UpdateTime = time.Now()
		//article.DeleteTime = "1970-01-01 00:00:00"
		id, err := ser.Dao.CreateAppArticle(article)
		if err != nil {
			return fullPath, imgName, fileName, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		req.ArticleId = id
	}
	_, imgPath, imgName, err := NewApiImgsService().SaveImage(req)
	if err != nil {
		return fullPath, imgName, fileName, err
	}
	fullPath, _ = url.JoinPath(imgPath, imgName)
	return fullPath, imgName, fileName, nil
}

func (ser *apiArticleService) SaveArticleDraft(req models.AppArticle) (articleDraft models.AppArticle, err error) {
	if req.AuthorId <= 0 {
		return articleDraft, errors.New(constant.ARTICLE_AUTHOR_ERR)
	}
	if req.Title != "" {
		articleDraft.Title = req.Title
	}
	if req.Content != "" {
		articleDraft.Content = req.Content
	}
	if req.Description != "" {
		articleDraft.Description = req.Description
	}
	if req.CoverUrl != "" {
		articleDraft.CoverUrl = req.CoverUrl
	}
	articleDraft.State = 1
	articleDraft.UpdateTime = time.Now()
	if req.Id != 0 {
		if req.Id <= 0 {
			return articleDraft, errors.New(constant.ARTICLE_UPDATE_ERR)
		}
		ser.Dao.UpdateArticle(req.Id, articleDraft)
	} else {
		articleDraft.CreateTime = time.Now()
		id, err := ser.Dao.CreateAppArticle(req)
		if err != nil {
			return articleDraft, errors.New(constant.ARTICLE_SAVE_ERR)
		}
		articleDraft.Id = id
	}
	return articleDraft, nil
}

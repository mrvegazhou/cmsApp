package api

import (
	"cmsApp/internal/constant"
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"cmsApp/pkg/utils/arrayx"
	"errors"
	"sync"
)

type apiArticleToolBarService struct {
	ArticleLikeDao      *dao.AppArticleLikeDao
	ArticleFavoritesDao *dao.AppArticleFavoritesDao
	FavoritesDao        *dao.AppFavoritesDao
}

var (
	instanceApiArticleToolBarService *apiArticleToolBarService
	onceApiArticleToolBarService     sync.Once
)

func NewApiArticleToolBarService() *apiArticleToolBarService {
	onceApiArticleToolBarService.Do(func() {
		instanceApiArticleToolBarService = &apiArticleToolBarService{
			ArticleLikeDao:      dao.NewAppArticleLikeDao(),
			ArticleFavoritesDao: dao.NewAppArticleFavoritesDao(),
			FavoritesDao:        dao.NewAppFavoritesDao(),
		}
	})
	return instanceApiArticleToolBarService
}

func (ser *apiArticleToolBarService) DoArticleLike(articleId, userId uint64) error {
	err := ser.ArticleLikeDao.DoArticleLike(articleId, userId)
	if err != nil {
		return errors.New(constant.ARTICLE_LIKE_ERR)
	}
	return nil
}

func (ser *apiArticleToolBarService) DoArticleUnlike(articleId, userId uint64) error {
	err := ser.ArticleLikeDao.DoArticleUnlike(articleId, userId)
	if err != nil {
		return errors.New(constant.ARTICLE_UNLIKE_ERR)
	}
	return nil
}

// 是否点赞，是否收藏，收藏夹列表
func (ser *apiArticleToolBarService) GetArticleToolBarData(articleId, userId uint64) (toolBarData models.AppArticleToolBarDataResp) {
	_, err := ser.ArticleLikeDao.CheckArticleLike(articleId, userId)
	if err != nil {
		toolBarData.IsLiked = false
	} else {
		toolBarData.IsLiked = true
	}
	favorites, err := ser.FavoritesDao.GetFavoritesByUser(userId)

	if err != nil {
		toolBarData.IsCollected = false
		toolBarData.Favorites = map[uint64]models.AppFavoritesItem{}
	} else {
		var favoritesIds []uint64
		favHasCheckedMap := make(map[uint64]models.AppFavoritesItem)
		for i := 0; i < len(favorites); i++ {
			favoritesIds = append(favoritesIds, favorites[i].Id)
			favHasChecked := models.AppFavoritesItem{}
			favHasChecked.IsChecked = false
			favHasChecked.Id = favorites[i].Id
			favHasChecked.Name = favorites[i].Name
			favHasCheckedMap[favorites[i].Id] = favHasChecked
		}
		articleFavorites := ser.ArticleFavoritesDao.GetArticleFavoritesList(favoritesIds)
		if len(articleFavorites) == 0 {
			toolBarData.IsCollected = false
			toolBarData.Favorites = favHasCheckedMap
		} else {
			toolBarData.IsCollected = true
			for i := 0; i < len(articleFavorites); i++ {
				tmp := favHasCheckedMap[articleFavorites[i].FavoritesId]
				if arrayx.IsContain(favoritesIds, articleFavorites[i].FavoritesId) {
					tmp.IsChecked = true
				}
				favHasCheckedMap[articleFavorites[i].FavoritesId] = tmp
			}
			toolBarData.Favorites = favHasCheckedMap
		}
	}
	return toolBarData
}

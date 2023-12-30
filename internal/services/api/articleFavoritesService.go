package api

import (
	"cmsApp/internal/dao"
	"sync"
)

type apiArticleFavoritesService struct {
	ArticleFavDao *dao.AppArticleFavoritesDao
	FavoritesDao  *dao.AppFavoritesDao
}

var (
	instanceApiArticleFavoritesService *apiArticleFavoritesService
	onceApiArticleFavoritesService     sync.Once
)

func NewApiArticleFavoritesService() *apiArticleFavoritesService {
	onceApiArticleFavoritesService.Do(func() {
		instanceApiArticleFavoritesService = &apiArticleFavoritesService{
			ArticleFavDao: dao.NewAppArticleFavoritesDao(),
			FavoritesDao:  dao.NewAppFavoritesDao(),
		}
	})
	return instanceApiArticleFavoritesService
}

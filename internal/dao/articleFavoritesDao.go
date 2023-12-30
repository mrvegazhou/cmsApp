package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppArticleFavoritesDao struct {
	DB *gorm.DB
}

var (
	instanceArticleFavorites *AppArticleFavoritesDao
	onceArticleFavoritesDao  sync.Once
)

func NewAppArticleFavoritesDao() *AppArticleFavoritesDao {
	onceArticleFavoritesDao.Do(func() {
		instanceArticleFavorites = &AppArticleFavoritesDao{DB: postgresqlx.GetDB(&models.AppArticleFavorites{})}
	})
	return instanceArticleFavorites
}

func (dao *AppArticleFavoritesDao) GetArticleFavoritesList(favoritesIds []uint64) []models.AppArticleFavorites {
	fav := []models.AppArticleFavorites{}
	dao.DB.Where("favorites_id IN (?)", favoritesIds).Find(&fav)
	return fav
}

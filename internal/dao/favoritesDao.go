package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppFavoritesDao struct {
	DB *gorm.DB
}

var (
	instanceFavorites *AppFavoritesDao
	onceFavoritesDao  sync.Once
)

func NewAppFavoritesDao() *AppFavoritesDao {
	onceFavoritesDao.Do(func() {
		instanceFavorites = &AppFavoritesDao{DB: postgresqlx.GetDB(&models.AppFavorites{})}
	})
	return instanceFavorites
}

func (dao *AppFavoritesDao) GetFavoritesByUser(userId uint64) ([]models.AppFavorites, error) {
	var favs []models.AppFavorites
	err := dao.DB.Where("user_id = ?", userId).Find(&favs)
	if err != nil {
		return favs, nil
	}
	return favs, nil
}

package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type AppReportDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	instanceAppReport *AppReportDao
	onceAppReportDao  sync.Once
)

func NewAppReportDao() *AppReportDao {
	onceAppReportDao.Do(func() {
		instanceAppReport = &AppReportDao{DB: postgresqlx.GetDB(&models.AppReport{})}
	})
	return instanceAppReport
}

func (dao *AppReportDao) CreateAppReport(report models.AppReport) (uint64, error) {
	if err := dao.DB.Create(&report).Error; err != nil {
		return 0, err
	}
	return report.Id, nil
}

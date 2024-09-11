package dao

import (
	"cmsApp/internal/models"
	"cmsApp/pkg/postgresqlx"
	"gorm.io/gorm"
	"sync"
)

type ReportReasonDao struct {
	DB *gorm.DB
	BaseDao
}

var (
	insReportReason     *ReportReasonDao
	onceReportReasonDao sync.Once
)

func NewReportReasonDao() *ReportReasonDao {
	onceReportReasonDao.Do(func() {
		insReportReason = &ReportReasonDao{DB: postgresqlx.GetDB(&models.ReportReason{})}
	})
	return insReportReason
}

func (dao *ReportReasonDao) GetReportReasonList(conditions map[string][]interface{}) ([]models.ReportReason, error) {
	reportReasons := []models.ReportReason{}
	Db := dao.DB
	if len(conditions) > 0 {
		Db = dao.ConditionWhere(Db, conditions, models.ReportReasonFields{}).Scopes(dao.Order("id desc"))
		if err := Db.Find(&reportReasons).Error; err != nil {
			return reportReasons, err
		}
	}
	return reportReasons, nil
}

func (dao *ReportReasonDao) GetAllReportReasonList() (reportReasons []models.ReportReason, err error) {
	reportReasons = []models.ReportReason{}
	err = dao.DB.Order("id desc").Find(&reportReasons).Error
	return
}

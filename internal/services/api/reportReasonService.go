package api

import (
	"cmsApp/internal/dao"
	"cmsApp/internal/models"
	"encoding/json"
	"sort"
	"sync"
)

type apiReportReasonService struct {
	Dao *dao.ReportReasonDao
}

var (
	instApiReportReasonService *apiReportReasonService
	onceApiReportReasonService sync.Once
)

func NewApiReportReasonService() *apiReportReasonService {
	onceApiReportReasonService.Do(func() {
		instApiReportReasonService = &apiReportReasonService{
			Dao: dao.NewReportReasonDao(),
		}
	})
	return instApiReportReasonService
}

func (ser *apiReportReasonService) GetReportReasons() (res []models.ReportReasonResp, err error) {
	allReasons, err := ser.Dao.GetAllReportReasonList()
	if err == nil {
		for _, reason := range allReasons {
			if reason.Pid == 0 {
				pNode := models.ReportReasonResp{}
				pNode.Id = reason.Id
				pNode.Name = reason.Name
				condition := models.ReportReasonCondition{}
				err = json.Unmarshal([]byte(reason.Conditions), &condition)
				if err == nil {
					pNode.Condition = condition
				}
				for _, child := range allReasons {
					if child.Pid == reason.Id {
						node := models.ReportReasonResp{}
						node.Id = child.Id
						node.Name = child.Name
						nodeCondition := models.ReportReasonCondition{}
						err = json.Unmarshal([]byte(child.Conditions), &nodeCondition)
						if err == nil {
							node.Condition = nodeCondition
						}
						pNode.Nodes = append(pNode.Nodes, node)
					}
				}
				sort.Slice(pNode.Nodes, func(i, j int) bool {
					return pNode.Nodes[i].Id < pNode.Nodes[j].Id
				})
				res = append(res, pNode)
			}
		}
		sort.Slice(res, func(i, j int) bool {
			return res[i].Id < res[j].Id
		})
	}
	return
}

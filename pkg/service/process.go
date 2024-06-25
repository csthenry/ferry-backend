package service

import (
	"encoding/json"
	"errors"
	"ferry/global/orm"
	"ferry/models/process"
	"ferry/models/system"
	"ferry/tools"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

type CreatorInfo struct {
	Name       string `json:"name"`
	Sex        string `json:"sex"`
	Mail       string `json:"mail"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

type WorkOrderData struct {
	process.WorkOrderInfo
	CurrentState string      `json:"current_state"`
	CreatorInfo  CreatorInfo `json:"creator_info"`
}

func ProcessStructure(c *gin.Context, processId int, workOrderId int) (result map[string]interface{}, err error) {
	var (
		processValue            process.Info
		processStructureDetails map[string]interface{}
		processNode             []map[string]interface{}
		tplDetails              []*process.TplInfo
		workOrderInfo           WorkOrderData
		workOrderTpls           []*process.TplData
		workOrderHistory        []*process.CirculationHistory
		stateList               []map[string]interface{}
		userInfo                system.SysUser
		roleInfo                system.SysRole
		deptInfo                system.Dept
		postInfo                system.Post
	)

	err = orm.Eloquent.Model(&processValue).Where("id = ?", processId).Find(&processValue).Error
	//if err != nil {
	//	err = fmt.Errorf("查询流程失败，%v", err.Error())
	//	return
	//}

	if processValue.Structure != nil && len(processValue.Structure) > 0 {
		err = json.Unmarshal([]byte(processValue.Structure), &processStructureDetails)
		if err != nil {
			err = fmt.Errorf("json转map失败，%v", err.Error())
			return
		}

		// 排序，使用冒泡
		p := processStructureDetails["nodes"].([]interface{})
		if len(p) > 1 {
			for i := 0; i < len(p); i++ {
				for j := 1; j < len(p)-i; j++ {
					if p[j].(map[string]interface{})["sort"] == nil || p[j-1].(map[string]interface{})["sort"] == nil {
						return nil, errors.New("流程未定义顺序属性，请确认")
					}
					leftInt, _ := strconv.Atoi(p[j].(map[string]interface{})["sort"].(string))
					rightInt, _ := strconv.Atoi(p[j-1].(map[string]interface{})["sort"].(string))
					if leftInt < rightInt {
						//交换
						p[j], p[j-1] = p[j-1], p[j]
					}
				}
			}
			for _, node := range processStructureDetails["nodes"].([]interface{}) {
				processNode = append(processNode, node.(map[string]interface{}))
			}
		} else {
			processNode = processStructureDetails["nodes"].([]map[string]interface{})
		}
	}

	processValue.Structure = nil
	result = map[string]interface{}{
		"process": processValue,
		"nodes":   processNode,
		"edges":   processStructureDetails["edges"],
	}

	// 获取历史记录
	err = orm.Eloquent.Model(&process.CirculationHistory{}).
		Where("work_order = ?", workOrderId).
		Order("id desc").
		Find(&workOrderHistory).Error
	if err != nil {
		return
	}
	result["circulationHistory"] = workOrderHistory

	if workOrderId == 0 {
		// 查询流程模版
		var tplIdList []int
		err = json.Unmarshal(processValue.Tpls, &tplIdList)
		if err != nil {
			err = fmt.Errorf("json转map失败，%v", err.Error())
			return
		}
		err = orm.Eloquent.Model(&tplDetails).
			Where("id in (?)", tplIdList).
			Find(&tplDetails).Error
		if err != nil {
			err = fmt.Errorf("查询模版失败，%v", err.Error())
			return
		}
		result["tpls"] = tplDetails
	} else {
		// 查询工单信息
		err = orm.Eloquent.Model(&process.WorkOrderInfo{}).
			Where("id = ?", workOrderId).
			Scan(&workOrderInfo).Error
		if err != nil {
			return
		}
		// 获取当前节点
		err = json.Unmarshal(workOrderInfo.State, &stateList)
		if err != nil {
			err = fmt.Errorf("序列化节点列表失败，%v", err.Error())
			return
		}
		if len(stateList) == 0 {
			err = errors.New("当前工单没有下一节点数据")
			return
		}

		// 整理需要并行处理的数据
		if len(stateList) > 1 {
		continueHistoryTag:
			for _, v := range workOrderHistory {
				status := false
				for i, s := range stateList {
					if v.Source == s["id"].(string) && v.Target != "" {
						status = true
						stateList = append(stateList[:i], stateList[i+1:]...)
						continue continueHistoryTag
					}
				}
				if !status {
					break
				}
			}
		}

		if len(stateList) > 0 {
		breakStateTag:
			for _, stateValue := range stateList {
				if processStructureDetails["nodes"] != nil {
					for _, processNodeValue := range processStructureDetails["nodes"].([]interface{}) {
						if stateValue["id"].(string) == processNodeValue.(map[string]interface{})["id"] {
							if _, ok := stateValue["processor"]; ok {
								for _, userId := range stateValue["processor"].([]interface{}) {
									if int(userId.(float64)) == tools.GetUserId(c) {
										workOrderInfo.CurrentState = stateValue["id"].(string)
										break breakStateTag
									}
								}
							} else {
								err = errors.New("未查询到对应的处理人字段，请确认。")
								return
							}
						}
					}
				}
			}

			if workOrderInfo.CurrentState == "" {
				workOrderInfo.CurrentState = stateList[0]["id"].(string)
			}
		}

		// 查询创建人信息
		err = orm.Eloquent.Model(&userInfo).Where("user_id = ?", workOrderInfo.Creator).Find(&userInfo).Error
		if err != nil {
			return
		}
		err = orm.Eloquent.Model(&deptInfo).Where("dept_id = ?", userInfo.DeptId).Find(&deptInfo).Error
		if err != nil {
			return
		}
		err = orm.Eloquent.Model(&postInfo).Where("post_id = ?", userInfo.PostId).Find(&postInfo).Error
		if err != nil {
			return
		}
		err = orm.Eloquent.Model(&roleInfo).Where("role_id = ?", userInfo.RoleId).Find(&roleInfo).Error
		if err != nil {
			return
		}
		workOrderInfo.CreatorInfo.Name = userInfo.NickName
		workOrderInfo.CreatorInfo.Sex = userInfo.Sex
		workOrderInfo.CreatorInfo.Mail = userInfo.Email
		workOrderInfo.CreatorInfo.Phone = userInfo.Phone
		workOrderInfo.CreatorInfo.Avatar = userInfo.Avatar
		workOrderInfo.CreatorInfo.Role = roleInfo.RoleName
		workOrderInfo.CreatorInfo.Department = deptInfo.DeptName
		workOrderInfo.CreatorInfo.Position = postInfo.PostName
		result["workOrder"] = workOrderInfo

		// 查询工单表单数据
		err = orm.Eloquent.Model(&workOrderTpls).
			Where("work_order = ?", workOrderId).
			Find(&workOrderTpls).Error
		if err != nil {
			return nil, err
		}
		result["tpls"] = workOrderTpls
	}
	return result, nil
}

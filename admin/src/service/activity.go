package service

import (
	"blitzSeckill/admin/src/db"
	"blitzSeckill/admin/src/model"
	"log"
)

type ActivityService struct {
}

func NewActivityService() *ActivityService {
	return &ActivityService{}
}

func (p *ActivityService) GetActivityList() ([]map[string]interface{}, error) {
	//_, err := model.ActivityNsp.GetActivityList(db.DB, 1,1)
	//if err != nil {
	//	log.Printf("ActivityEntity.GetActivityList, err : %v", err)
	//	return nil, err
	//}

	//for _, v := range activityList {
	//	startTime, _ := com.StrTo(fmt.Sprint(v["start_time"])).Int64()
	//	v["start_time_str"] = time.Unix(startTime, 0).Format("2006-01-02 15:04:05")
	//
	//	endTime, _ := com.StrTo(fmt.Sprint(v["end_time"])).Int64()
	//	v["end_time_str"] = time.Unix(endTime, 0).Format("2006-01-02 15:04:05")
	//
	//	nowTime := time.Now().Unix()
	//	if nowTime > endTime {
	//		v["status_str"] = "已结束"
	//		continue
	//	}
	//
	//	status, _ := com.StrTo(fmt.Sprint(v["status"])).Int()
	//	if status == model.ActivityStatusNormal {
	//		v["status_str"] = "正常"
	//	} else if status == model.ActivityStatusDisable {
	//		v["status_str"] = "已禁用"
	//	}
	//}
	//
	//log.Printf("get activity success, activity list is [%v]", activityList)
	return []map[string]interface{}{}, nil
}

func (p *ActivityService) CreateActivity(activity *model.Activity) error {
	//写入到数据库
	err := model.ActivityNsp.Create(db.DB, activity)
	if err != nil {
		log.Printf("ActivityModel.CreateActivity, err : %v", err)
		return err
	}
	return nil
}

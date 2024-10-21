package logic

import (
	"GinBlog/dao/mysql"
	"GinBlog/models"
	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	//查找到所有的community并且返回
	data, err = mysql.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		return nil, err
	}
	return data, err
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}

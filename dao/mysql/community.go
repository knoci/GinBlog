package mysql

import (
	"GinBlog/models"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

func GetCommunityList() (data []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&data, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (detail *models.CommunityDetail, err error) {
	detail = new(models.CommunityDetail)
	sqlStr := "select community_id, community_name, introduction, create_time from community where community_id = ?"
	if err := db.Get(detail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidID
		}
	}
	fmt.Println("%v", detail)
	return detail, err
}

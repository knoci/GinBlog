package mysql

import (
	"GinBlog/models"
	"GinBlog/setting"
	"testing"
)

func init() {
	dbCfg := setting.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "123",
		DbName:       "ginblog",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		AuthorID:    123,
		CommunityID: 1,
		Title:       "test",
		Content:     "just a simple test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}

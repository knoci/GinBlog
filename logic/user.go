package logic

import (
	"GinBlog/dao/mysql"
	"GinBlog/models"
	"GinBlog/pkg/jwt"
	"GinBlog/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 合法性判断
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 保存到数据库
	user := models.User{
		UserID:   snowflake.GenID(),
		Username: p.Username,
		Password: p.Password,
	}
	err = mysql.InsertUser(&user)
	if err != nil {
		return err
	}
	return nil
}

func Login(p *models.ParamLogin) (*models.User,error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 传递的是指针，就能拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return user, err
}

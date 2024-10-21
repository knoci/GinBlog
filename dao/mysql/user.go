package mysql

import (
	"GinBlog/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const salt = "knoci1337"


func InsertUser(user *models.User) (err error) {
	// 密码加密
	user.Password = encryptPassword(user.Password)
	sqlStr := `insert into user(user_id, username, password) value(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	if err != nil {
		return
	}
	return nil
}

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrUserExist
	}
	return nil
}

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id, username , password from user where username = ?`
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrUserNotExist
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrInvalidPassword
	}
	return
}

// GetUserById 根据id获取用户信息
func GetUserById(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, uid)
	return
}
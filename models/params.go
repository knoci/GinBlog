package models

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	PostID string `json:"post_id" binding:"required"`
	// 赞成1反对-1取消投票0,validator的oneof限定
	Direction int8 `json:"direction,string" binding:"oneof=1 0 -1" `
}

type ParamPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page int64 `form:"page"`
	Size int64 `form:"size"`
	Order string `form:"order"`
}

type ParamCommunityPostList struct {

}
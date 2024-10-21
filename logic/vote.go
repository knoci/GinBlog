package logic

import (
	"GinBlog/dao/redis"
	"GinBlog/models"
	"go.uber.org/zap"
	"strconv"
)



func PostVote(userID int64, p *models.ParamVoteData) error{
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}

package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"math"
	"time"
)

const (
	oneWeek = 7 * 24 * 3600
	scorePerVote = 2880	// 一票加2880，2880*30 = 86400 一天
)

var (
	ErrVoteTimeExpire = errors.New("投票过时")
	ErrVoteRepeat = errors.New("重复投票")
)


/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票    --> 更新分数和投票记录 差值的绝对值是1 +2880
	2. 之前投反对票，现在改投赞成票    --> 更新分数和投票记录 差值的绝对值是2 + 2880*2
direction=0时，有两种情况：
	1. 之前投过赞成票，现在要取消投票  --> 更新分数和投票记录 差值的绝对值是1 -2880
	2. 之前投过反对票，现在要取消投票  --> 更新分数和投票记录 差值的绝对值是1 +2880
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票    --> 更新分数和投票记录 差值的绝对值是1 -2880
	2. 之前投赞成票，现在改投反对票    --> 更新分数和投票记录 差值的绝对值是2 -2880*2

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 到期之后删除那个 KeyPostVotedZSetPF
*/
func VoteForPost(userID, postID string, value float64) error {
	// 判断投票情况
	ctx := context.Background()
	postTime := client.ZScore(ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix()) - postTime > oneWeek {
		return ErrVoteTimeExpire
	}
	// 更新分数

	pipeline := client.TxPipeline()
	// 记录用户投票
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		oldVal := client.ZScore(ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val() // 查询投票记录
		var op float64
		if value == oldVal {
			return ErrVoteRepeat
		}
		if value > oldVal {
			op = 1
		} else {
			op = -1
		}
		diff := math.Abs(oldVal - value)
		pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZSetPF+postID), &redis.Z{
			Score: value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}

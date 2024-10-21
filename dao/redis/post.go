package redis

import (
	"GinBlog/models"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	ctx := context.Background()
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	// 更新：把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)
	return err
}

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 确定所有起点
	start := (page -1) * size
	end := start + size -1
	// 查询,按分数从大到小
	ctx := context.TODO()
	return  client.ZRevRange(ctx, key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 获取排序顺序
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

// 根据ids查询每篇帖子的投票数据
func GetPostVoteData(ids []string) (data []int64, err error){
	ctx := context.TODO()
	data = make([]int64, 0, len(ids))
	// 使用pipeline一次发送剁掉指令减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF+id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// 按社区根据ids查询每篇帖子的票
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OderScore  {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用 zinterstore 把分区的帖子set与帖子分数的 zset 生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据
	// 社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存key减少zinterstore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	ctx := context.Background()
	if client.Exists(ctx, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:    []string{cKey, orderKey},
			Aggregate: "MAX",
		}) // zinterstore 计算
		pipeline.Expire(ctx, key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 存在的话就直接根据key查询ids
	return getIDsFormKey(key, p.Page, p.Size)
}
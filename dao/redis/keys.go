package redis

// redis key注意使用命名空间方式方便查询和拆分
const (
	Prefix = "ginblog:"
	KeyPostTimeZSet = "post:time" // Zset 发帖时间为分数
	KeyPostScoreZSet = "post:score" // Zset 帖子评分为分数
	KeyPostVotedZSetPF = "post:voted" // zset 记录用户及投票类型，参数是post id
	KeyCommunitySetPF = "community:" // set 保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}
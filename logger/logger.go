package logger

import (
	"GinBlog/setting"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack" // lumberjack 是一个简单的日志文件滚动库
	"go.uber.org/zap"                 // zap 是一个快速的、结构化的、可靠的 Go 日志库
	"go.uber.org/zap/zapcore"         // zap 的核心包
)

// lg 是一个全局的 zap.Logger 实例
var lg *zap.Logger

// Init 初始化日志配置
func Init(cfg *setting.LogConfig, mode string) (err error) {
	// 获取日志写入器
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	// 获取编码器
	encoder := getEncoder()

	// 解析日志级别
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return // 如果解析失败，返回错误
	}

	// 创建 zap 核心对象
	var core zapcore.Core
	if mode == "dev" {
		// 如果是开发模式，日志同时输出到终端和文件
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee( // Tee 表示同时写入多个 Writer
			//zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// 生产模式，只输出到文件
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	// 创建 zap 日志对象
	lg = zap.New(core, zap.AddCaller()) // 添加调用者信息
	// 替换全局日志对象
	zap.ReplaceGlobals(lg)
	// 记录初始化日志
	zap.L().Info("init logger success")
	return
}

// getEncoder 创建并配置日志编码器
func getEncoder() zapcore.Encoder {
	// 使用生产环境的编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间编码器
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 设置时间字段名称
	encoderConfig.TimeKey = "time"
	// 设置日志级别编码器
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 设置持续时间编码器
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// 设置调用者编码器
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 创建 JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 创建并配置日志写入器
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 使用 lumberjack 作为日志文件的写入器
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 日志文件路径
		MaxSize:    maxSize,   // 文件最大大小
		MaxBackups: maxBackup, // 最多备份文件数量
		MaxAge:     maxAge,    // 文件最长保存天数
	}
	// 将 lumberjack 写入器包装为 zap 写入器
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 是一个 Gin 中间件，用于记录 HTTP 请求日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		// 获取请求路径和查询字符串
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		// 继续处理请求
		c.Next()

		// 计算请求处理时间
		cost := time.Since(start)
		// 记录日志
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery 是一个 Gin 中间件，用于捕获和记录 panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用 defer 延迟执行 panic 恢复
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是连接断开导致的错误
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 记录请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					// 如果是连接断开，记录错误日志
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接已断开，记录错误并终止请求
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				// 如果需要记录 stack trace
				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 设置 HTTP 状态码并终止请求
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 继续处理请求
		c.Next()
	}
}

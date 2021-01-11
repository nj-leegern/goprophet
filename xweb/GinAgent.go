package xweb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/nj-leegern/goprophet/utils"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

/*
	Gin agent
*/

type GinAgent struct {
	engine *gin.Engine
	port   uint64
}

// 路由配置项
type Options struct {
	Method string                     // GET POST PUT OPTIONS
	Routes map[string]gin.HandlerFunc // path -> handle
}

/* 创建Gin agent实例 */
func NewGinAgent(port uint64, logPath string, middleware []gin.HandlerFunc, options ...Options) *GinAgent {

	fmt.Println("initializing gin agent ...")

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	// default middleware
	engine.Use(generateGinLog(logPath), cors())
	// user middleware
	if middleware != nil && len(middleware) > 0 {
		engine.Use(middleware...)
	}

	// routes
	for _, option := range options {
		routes := option.Routes
		if strings.ToUpper(option.Method) == "GET" {
			for path, handle := range routes {
				if !strings.HasPrefix(path, "/") {
					path = "/" + path
				}
				engine.GET(path, handle)
			}
		}
		if strings.ToUpper(option.Method) == "POST" {
			for path, handle := range routes {
				if !strings.HasPrefix(path, "/") {
					path = "/" + path
				}
				engine.POST(path, handle)
			}
		}
		if strings.ToUpper(option.Method) == "PUT" {
			for path, handle := range routes {
				if !strings.HasPrefix(path, "/") {
					path = "/" + path
				}
				engine.PUT(path, handle)
			}
		}
		if strings.ToUpper(option.Method) == "OPTIONS" {
			for path, handle := range routes {
				if !strings.HasPrefix(path, "/") {
					path = "/" + path
				}
				engine.OPTIONS(path, handle)
			}
		}
	}

	return &GinAgent{engine: engine, port: port}
}

/* 启动Gin agent */
func (ga *GinAgent) RunAgent() error {
	port := strconv.FormatUint(ga.port, 10)
	fmt.Printf("gin agent start with port %s\n", port)
	err := ga.engine.Run(":" + port)
	return err
}

// Gin默认输出日志
func generateGinLog(logPath string) gin.HandlerFunc {
	logger := log.New()
	// 目录不存在则创建
	if !utils.IsExist(logPath) {
		os.MkdirAll(logPath, os.ModePerm)
	}
	fileName := path.Join(logPath, "gin-api.log")
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil {
		fmt.Println("create gin log failed", err)
	}
	// 设置日志输出的路径
	logger.Out = src
	logger.SetLevel(log.DebugLevel)
	logger.SetReportCaller(true)
	logWriter, err := rotatelogs.New(
		fileName+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(fileName),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	writeMap := lfshook.WriterMap{
		log.InfoLevel:  logWriter,
		log.FatalLevel: logWriter,
		log.DebugLevel: logWriter, // 为不同级别设置不同的输出目的
		log.WarnLevel:  logWriter,
		log.ErrorLevel: logWriter,
		log.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: utils.DATE_PATTERN_DEFAULT,
	})
	logger.AddHook(lfHook)

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		// 日志格式: 状态码、执行时间、请求ip、请求方法、请求路由
		logger.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

// 请求跨域
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}

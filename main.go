package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration structure
type Config struct {
	Basicauth struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"basicauth"`
	Logbasepath string `yaml:"logbasepath"`
	Port        string `yaml:"port"`
}

// LogReq represents the request payload for logging
type LogReq struct {
	P   string `json:"p" xml:"p" form:"p"`
	Msg string `json:"msg" xml:"msg" form:"msg"`
}

func main() {
	// 读取配置文件
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file: %v", err)
	}

	fmt.Println(config)
	// 创建 Fiber 应用
	app := fiber.New()

	base := app.Group("/api")

	// Basic Auth 中间件配置
	v1 := base.Group("/v1", basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.Basicauth.Username: config.Basicauth.Password,
		},
	}))

	// 定义中间件处理函数
	v1.Use(func(c *fiber.Ctx) error {
		req := new(LogReq)
		// 通过 BodyParser 将 body 参数绑定到变量 req
		if err := c.BodyParser(req); err != nil {
			return err
		}

		// 获取参数
		logFilePath := req.P
		logMessage := req.Msg

		// 处理缺少参数情况
		if logFilePath == "" || logMessage == "" {
			date := time.Now().Format("2006-01-02")
			logFilePath = fmt.Sprintf("/tmp/error-%s.log", date)
		}

		logFilePath = config.Logbasepath + logFilePath
		// 检查并创建日志文件目录
		logFileDirPath := path.Dir(logFilePath)
		if _, err := os.Stat(logFileDirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(logFileDirPath, os.ModePerm); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to create log directory: %v", err))
			}
		}

		// 打开或创建日志文件
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to open log file: %v", err))
		}
		defer logFile.Close()

		// 配置日志中间件
		customLogger := logger.New(logger.Config{
			Output:     logFile,
			Format:     "[${time}] ${method} ${protocol} ${host} ${url} ${status}  ${latency} |${ua}|${reqHeaders}|${queryParams}|${body}|${resBody}\n",
			TimeFormat: "2006-01-02 15:04:05",
			TimeZone:   "Asia/Shanghai",
		})

		// 使用自定义日志中间件处理当前请求
		return customLogger(c)
	})

	// /v1/log 路由处理
	v1.Post("/log", func(c *fiber.Ctx) error {
		req := new(LogReq)
		// 通过 BodyParser 将 body 参数绑定到变量 req
		if err := c.BodyParser(req); err != nil {
			return err
		}

		// 检查参数是否缺失
		if req.P == "" || req.Msg == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing 'p' or 'msg' parameter")
		}

		// 返回成功响应
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": true,
			"msg":    "ok",
		})
	})

	// 启动服务器
	log.Fatal(app.Listen(":" + config.Port))
}

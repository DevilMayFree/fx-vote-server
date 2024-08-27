package initialize

import (
	"context"
	"errors"
	"fmt"
	"fx-vote-server/common/app"
	"fx-vote-server/common/constant"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type SetVoteRequestData struct {
	List []string `json:"list"`
}

type VoteRequestData struct {
	Num string `json:"num"`
}

type PageData struct {
	Num   int    `json:"num"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// 创建全局互斥锁
var mu sync.Mutex

// 判断请求是否是接口请求
func isAPIRequest(path string) bool {
	if len(path) < 5 {
		return false
	}
	if len(path) > 5 && path[0] == '/' && path[1:5] == "api/" {
		return true
	}
	return false
}

// cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func runServer() {
	// init gin server

	routers := Routers()

	routers.LoadHTMLGlob("templates/*")

	// 模拟的用户数据
	var users = map[string]string{
		"admin": "123456", // username: password
	}

	// 跨域设置
	routers.Use(Cors())

	// 认证中间件
	routers.Use(func(c *gin.Context) {
		// 检查请求路径是否以 "/api/" 开头
		if !isAPIRequest(c.Request.URL.Path) {
			cookie, err := c.Cookie("fx-vote-server")
			if err != nil || cookie != "authenticated" {
				c.HTML(http.StatusUnauthorized, "login.html", nil)
				c.Abort()
				return
			}
		}
		c.Next()
	})

	// 登录页面
	routers.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// 首页处理
	routers.GET("/index", func(c *gin.Context) {
		list := app.Redis.LRange(context.Background(), constant.RedisVoteKey, 0, 7)
		val := list.Val()

		var dataList []PageData
		for i, item := range val {
			name := constant.TeacherList[i]
			d := PageData{
				Num:   i,
				Name:  name,
				Value: item,
			}
			dataList = append(dataList, d)
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"dataList": dataList,
		})
	})

	api := routers.Group("/api")
	{
		// 登录处理
		api.POST("/login", func(c *gin.Context) {
			username := c.PostForm("username")
			password := c.PostForm("password")

			trimUserName := strings.TrimSpace(username)
			trimPassword := strings.TrimSpace(password)

			// 验证用户
			if pass, exists := users[trimUserName]; exists && pass == trimPassword {
				// 设置会话或 token 这里使用 Cookie 示例
				c.SetCookie("fx-vote-server", "authenticated", 3600, "/", "", false, true)
				c.Redirect(http.StatusSeeOther, "/index")
			} else {
				// c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
				c.HTML(http.StatusOK, "login.html", gin.H{
					"ErrorMessage": "Invalid username or password",
				})
			}
		})

		// 投票api
		api.POST("/vote", func(c *gin.Context) {

			var data VoteRequestData
			if err := c.BindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
				return
			}

			// 客户端ip
			ip := c.ClientIP()

			// 检查 IP 地址是否存在
			exists, existsErr := app.Redis.SIsMember(context.Background(), constant.RedisIpKey, ip).Result()
			if existsErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"result": "ip_exist"})
				return
			}

			if exists {
				c.JSON(http.StatusBadRequest, gin.H{"result": "IP address exists"})
				return
			}

			// 如果 IP 地址不存在，将其添加到 Redis
			_, addErr := app.Redis.SAdd(context.Background(), constant.RedisIpKey, ip).Result()
			if addErr != nil {
				f, _ := fmt.Printf("Error adding IP to Redis: %v", addErr)
				c.JSON(http.StatusInternalServerError, gin.H{"result": f})
				return
			}

			_, expireErr := app.Redis.Expire(context.Background(), constant.RedisIpKey, constant.RedisIpExpiration).Result()
			if expireErr != nil {
				f, _ := fmt.Printf("Error adding IP to Redis: %v", expireErr)
				c.JSON(http.StatusInternalServerError, gin.H{"result": f})
				return
			}

			mu.Lock()
			defer mu.Unlock()

			// 获取指定索引的值
			index, _ := strconv.Atoi(data.Num)
			value, lIndexError := app.Redis.LIndex(context.Background(), constant.RedisVoteKey, int64(index)).Result()
			if lIndexError != nil {
				f, _ := fmt.Printf("Error getting value at index %d: %v", index, lIndexError)
				c.JSON(http.StatusInternalServerError, gin.H{"result": f})
				return
			}

			// 将值转换为整数并加1
			intValue, atoiError := strconv.Atoi(value)
			if atoiError != nil {
				f, _ := fmt.Printf("Error converting value to integer: %v", atoiError)
				c.JSON(http.StatusInternalServerError, gin.H{"result": f})
				return
			}
			intValue++

			// 将更新后的值转换回字符串
			newValue := strconv.Itoa(intValue)

			// 更新列表中的值
			_, lSetError := app.Redis.LSet(context.Background(), constant.RedisVoteKey, int64(index), newValue).Result()
			if lSetError != nil {
				f, _ := fmt.Printf("Error setting value at index %d: %v", index, lSetError)
				c.JSON(http.StatusInternalServerError, gin.H{"result": f})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"result": "success",
			})
		})

		// 获取各位当前投票数
		api.GET("/getVote", func(c *gin.Context) {
			list := app.Redis.LRange(context.Background(), constant.RedisVoteKey, 0, 7)
			val := list.Val()

			c.JSON(http.StatusOK, gin.H{
				"result": val,
			})
		})

		// 设置各位的投票数
		api.POST("/setVote", func(c *gin.Context) {

			var data SetVoteRequestData
			if err := c.BindJSON(&data); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
				return
			}

			clearErr := app.Redis.Del(context.Background(), constant.RedisVoteKey).Err()
			if clearErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"clear error": clearErr.Error()})
				return
			}

			voteList := data.List
			err := app.Redis.RPush(context.Background(), constant.RedisVoteKey, voteList).Err()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"set error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"result": "success",
			})
		})
	}

	p := app.Config.Server.Port
	port := ":" + strconv.Itoa(p)

	server := &http.Server{
		Addr:    port,
		Handler: routers,
	}

	startCron()
	app.GinLog.Info("server started listen" + server.Addr)
	gracefulShutdown(server)
}

// gracefulShutdown gin优雅关闭
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	server.RegisterOnShutdown(func() {
		app.GinLog.Info("start execute out shutdown")
	})

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				app.GinLog.Info("Server closed under request")
			} else {
				app.Log.Info("Server closed unexpect")
				os.Exit(1)
			}
		}
	}()

	<-quit
	app.GinLog.Info("receive closeServer signal")

	if err := server.Shutdown(context.Background()); err != nil {
		app.Log.Error("ServerClose:", zap.Error(err))
		os.Exit(1)
	}
	releaseRes()
	app.GinLog.Info("Server exiting")
	os.Exit(0)
}

// releaseRes 释放资源
func releaseRes() {

	// release Redis Resource
	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			app.Log.Error("ServerClose", zap.NamedError("close Redis error", err))
		} else {
			app.GinLog.Info("close Redis success")
		}
	}
}

// 定时清理
func startCron() {
	// 初始化定时任务调度器
	c := newWithSeconds()
	// 添加定时任务，每天 00:01 执行
	// _, err := c.AddFunc("0 01 00 ? * *", func() {
	_, err := c.AddFunc("0 0/3 * * * ?", func() {
		err := deleteKey()
		if err != nil {
			log.Printf("Failed to delete key: %v", err)
		} else {
			log.Println("Key deleted successfully.")
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// 启动定时任务调度器
	c.Start()
}

// 支持6位cron表达式
func newWithSeconds() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

// deleteKey 从 Redis 中删除指定的 key
func deleteKey() error {
	_, err := app.Redis.Del(context.Background(), constant.RedisIpKey).Result()
	if err != nil {
		return fmt.Errorf("could not delete key %s: %w", constant.RedisIpKey, err)
	}
	return nil
}

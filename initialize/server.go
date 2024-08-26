package initialize

import (
	"context"
	"errors"
	"fx-vote-server/common/app"
	"fx-vote-server/common/constant"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

type SetVoteRequestData struct {
	List []string `json:"list"`
}

func runServer() {
	// init gin server

	routers := Routers()

	routers.LoadHTMLGlob("templates/*")

	// 模拟的用户数据
	var users = map[string]string{
		"admin": "123456", // username: password
	}

	// 登录处理
	routers.POST("/login", func(c *gin.Context) {
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

	// 认证中间件
	routers.Use(func(c *gin.Context) {
		cookie, err := c.Cookie("fx-vote-server")
		if err != nil || cookie != "authenticated" {
			c.HTML(http.StatusUnauthorized, "login.html", nil)
			c.Abort()
			return
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

		c.HTML(http.StatusOK, "index.html", gin.H{
			"val": val,
		})
	})

	// 获取各位当前投票数
	routers.GET("/getVote", func(c *gin.Context) {
		list := app.Redis.LRange(context.Background(), constant.RedisVoteKey, 0, 7)
		val := list.Val()

		c.JSON(http.StatusOK, gin.H{
			"result": val,
		})
	})

	// 设置各位的投票数
	routers.POST("/setVote", func(c *gin.Context) {

		var data SetVoteRequestData
		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
			return
		}

		clearErr := app.Redis.Del(context.Background(), constant.RedisVoteKey).Err()
		if clearErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"clear error": clearErr.Error()})
		}

		voteList := data.List
		err := app.Redis.RPush(context.Background(), constant.RedisVoteKey, voteList).Err()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"set error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
	})

	/*routers.GET("/ping", func(c *gin.Context) {
		app.Redis.Set(context.Background(), "a", "b", time.Hour)
		// time.Sleep(10 * time.Second)
		get := app.Redis.Get(context.Background(), "a").String()
		app.Log.Info("get redis:", zap.String("get:", get))
		c.JSON(http.StatusOK, gin.H{
			"result": "pong",
		})
	})

	routers.GET("/err", func(c *gin.Context) {
		panic("oh error happen")
	})*/

	p := app.Config.Server.Port
	port := ":" + strconv.Itoa(p)

	server := &http.Server{
		Addr:    port,
		Handler: routers,
	}

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

package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-backend-test/core"
	"github.com/go-backend-test/infra"
)

type RestfulVar struct {
	Port string `mapstructure:"PORT"`
}

type Restful struct {
	user core.IUserCore
	line core.ILine
	vars RestfulVar
}

func NewRestful() *Restful {
	rf := &Restful{}
	rf.user = core.NewUserCore()
	initEnv := func() {
		if err := infra.GetEnv(&rf.vars); err != nil {
			log.Fatal(err)
		}
	}

	initEnv()

	rf.line = core.NewLineCore()

	return rf
}

func (rf *Restful) Run(wg *sync.WaitGroup) {
	r := gin.Default()

	// 添加一个路由处理函数，处理 GET 请求
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	r.POST("/send-message", rf.line.HandleSendMessageByAPI())

	r.POST("/callback", rf.line.HandleCallback())

	r.GET("/query-user-message/:user_id", func(ctx *gin.Context) {
		id := ctx.Param("user_id")

		msgs, err := rf.user.ListUserMessageByUserId(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "server internal error",
			})
			return
		}

		ctx.JSON(http.StatusOK, msgs)

	})

	if err := r.Run(rf.vars.Port); err != nil {
		panic(err)
	}

	wg.Done()
}

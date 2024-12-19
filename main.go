package main

import (
	"RR7/api"
	"RR7/dao"
	"RR7/flag"
	"RR7/models"
	"RR7/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func main() {
	models.DB = dao.InitMysql()

	option := flag.Parse()
	if flag.IsWebStop(option) {
		flag.SwitchOption(option)
		return
	}

	r := gin.Default()
	r.POST("register", api.Register)
	r.POST("login", api.Login)
	r.GET("/messages/like", api.GetMessageLikes)

	p := r.Group("/")
	p.Use(utils.JWTAuthMiddleware())

	p.POST("/messages", api.PostLog)
	p.GET("/messages", api.GetLogs)
	p.DELETE("/messages", api.DeleteLog)
	p.POST("/messages/like", api.LikeMessage)
	p.DELETE("messages/like", api.CancelLikeMessage)

	r.Run(":8080")
}

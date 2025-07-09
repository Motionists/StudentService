package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化数据库连接
	initDB()
	defer DB.Close()

	router := gin.Default()

	router.POST("/login", login)
	api := router.Group("/api/v1")
	api.Use(jwtMiddleware()) // 添加 JWT 验证中间件
	{
		api.GET("/students", ListStudents)
		api.POST("/students", CreateStudent)
		api.GET("/students/:id", GetStudent)
		api.PUT("/students/:id", UpdateStudent)
		api.DELETE("/students/:id", DeleteStudent)
	}
	log.Println("服务器启动在 :8080")
	router.Run(":8080")
}

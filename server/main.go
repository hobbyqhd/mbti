package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 连接数据库
	dbURL := os.Getenv("DATABASE_URL")
	var err error
	DB, err = sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	// 设置数据库日志
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(time.Hour)

	// 添加SQL语句日志记录
	if os.Getenv("LOG_LEVEL") == "debug" {
		DB.SetConnMaxLifetime(time.Hour)
		log.Println("SQL日志记录已启用")
	}

	// 设置CORS中间件
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API路由
	r.GET("/api/questions", GetQuestions)
	r.POST("/api/submit", SubmitAnswers)
	r.GET("/api/result/:id", GetResult)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}

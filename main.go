package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
  "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID uint `gorm:"primary;type: bigint AUTO_INCREMENT"`
	Title string `gorm:"type: varchar(100)"`
	Slug string `gorm:"unique;type: varchar(100)"`
	Desc string `gorm:"type: text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var DB *gorm.DB

func main() {
	// database setup
	var err error
	dsn := "didik27:Didik.,.88@tcp(127.0.0.1:3306)/learn_gin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	DB.AutoMigrate(&Article{})

	if err != nil {
    panic("failed to connect database")
  }

	router := gin.Default()

	APIv1 := router.Group("/api/v1")
	{
		APIv1.GET("/", getWelcome)
		article := APIv1.Group("/article")
		{
			article.GET("/", getArticles)
			article.GET("/:slug", getDetailArticle)
			article.POST("/", addArticle)
		}
	}

	router.Run(":3000")
}

func getArticles(c *gin.Context) {
	articles := []Article{}
	DB.Find(&articles)

	c.JSON(200, gin.H{
		"success": true,
		"message": "welcome to go gin framework",
		"data": articles,
	})
}

func getDetailArticle(c *gin.Context) {
	slug := c.Param("slug")
	
	article := Article{}
	err := DB.First(&article, "slug = ?", slug).Error
	
	if err == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"success": false,
			"message": "data not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "welcome to go gin framework",
		"data": article,
	})
}

func getWelcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"message": "welcome to go gin framework",	
	})
}

func addArticle(c *gin.Context) {
	article := Article {
		Title: c.PostForm("title"),
		Desc: c.PostForm("desc"),
		Slug: slug.Make(c.PostForm("title")),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	DB.Create(&article)

	c.JSON(201, gin.H{
		"success": true,
		"data": article,
	})
}
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
)

var db *gorm.DB

type Todo struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func InitMysql() (err error) {
	username := "root"
	password := "xwt011028"
	host := "localhost"
	part := "3306"
	DbName := "test"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, part, DbName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	return
}

func main() {
	err := InitMysql()
	if err != nil {
		panic(err)
	}
	//db.AutoMigrate(&Todo{})
	r := gin.Default()
	//指定static文件位置
	r.Static("/static", "static")
	//指定模板位置
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加一个待办事项
		v1Group.POST("/todo", func(c *gin.Context) {
			//1.拿到数据
			var todo Todo
			c.BindJSON(&todo)
			//2.存储数据
			err = db.Create(&todo).Error
			//3.返回响应
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//查看某个待办事项
		//查看全部待办事项
		v1Group.GET("/todo", func(c *gin.Context) {
			var todolist []Todo
			err = db.Find(&todolist).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, todolist)
			}
		})
		//修改某个待办事项
		v1Group.PUT(("/todo/:id"), func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": "id无效",
				})
			}
			var todo Todo
			err = db.Where("id = ?", id).First(&todo).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
			}
			c.BindJSON(&todo)
			db.Save(&todo)
			c.JSON(http.StatusOK, todo)
		})
		//删除某个待办事项
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": "id无效",
				})
			}
			err = db.Where("id = ?", id).Delete(&Todo{}).Error
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"err": err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					id: "删除成功",
				})
			}
		})
	}
	r.Run(":8080")
}

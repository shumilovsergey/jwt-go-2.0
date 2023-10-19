package main

import (
	"jwt/controllers"
	"jwt/initializers"
	"jwt/middleware"

	"github.com/gin-gonic/gin"
)

// действия до основной функции
func init() {
	initializers.LoadEnvVaribles()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

// точка входа
func main() {
	//запуск web сервера
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "MAIN PAGE",
		})
	})

	//controllers
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run() // listen and serve on 0.0.0.0:env
}

package main

import (
	"fmt"
	"root/src/config"
	"root/src/controllers"
	"root/src/core/db"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	fmt.Println("Im a dummy!")
	// Pass on to the next-in-chain
	c.Next()
}
func main() {
	r := gin.Default()

	//db.InitRedis(1)
	db.InitMongoDB()
	//db.InitPostgresDB()
	//db.InitGorm()

	//models.MigrateUsers()

	// register controllers
	controllers.AuthController(r)

	r.Use(authMiddleware)
	controllers.UsersController(r)
	controllers.ArticlesController(r)
	controllers.ProductsController(r)
	controllers.SwaggersController(r)

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// running
	r.Run(fmt.Sprintf(":%s", config.LoadConfig("PORT")))
}

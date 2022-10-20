package main

import (
	"orders/internal/db"
	"orders/internal/handler"
	rout "orders/internal/router"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	db := db.NewDB()
	handler := handler.InitHandler(db)
	rout.InitRoutes(router, handler)
	router.Run(":8080")
}

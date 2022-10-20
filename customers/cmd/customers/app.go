package main

import (
	"customers/internal/db"
	"customers/internal/handler"
	rout "customers/internal/router"

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

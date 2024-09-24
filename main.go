package main

import (
	"github.com/gin-gonic/gin"
	"wardenapi/internal/services"
)

func main() {
	router := gin.Default()
	services.MapRunePriceRoutes(router)

	err := router.Run(":3000")
	if err != nil {
		panic(err)
	}
}

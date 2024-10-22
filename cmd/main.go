package main

import (
	"log"

	"github.com/Ekreke/whatshappening/internal/wsh"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routeController, err := InitializeRouteController()
	if err != nil {
		log.Fatalf("Failed to initialize route controller: %v", err)
	}

	wsh.InstallRouters(r, routeController)

	r.Run(":8080")
}

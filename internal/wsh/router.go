package wsh

import (
	"github.com/Ekreke/whatshappening/internal/wsh/controller"
	"github.com/gin-gonic/gin"
)

func InstallRouters(g *gin.Engine, rc *controller.RouteController) {
	api := g.Group("/api")
	api.GET("/routes", rc.GetRoutes)
}

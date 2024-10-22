package controller

import (
	"net/http"

	"github.com/Ekreke/whatshappening/internal/wsh/biz"
	"github.com/gin-gonic/gin"
)

type RouteController struct {
	routeBiz *biz.RouteBiz
}

func NewRouteController(routeBiz *biz.RouteBiz) *RouteController {
	return &RouteController{
		routeBiz: routeBiz,
	}
}

func (rc *RouteController) GetRoutes(c *gin.Context) {
	response, err := rc.routeBiz.FetchAndFilterRoutes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

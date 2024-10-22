package main

import (
	"github.com/Ekreke/whatshappening/internal/wsh/biz"
	"github.com/Ekreke/whatshappening/internal/wsh/controller"
	"github.com/google/wire"
)

func InitializeRouteController() (*controller.RouteController, error) {
	wire.Build(
		biz.NewRouteBiz,
		controller.NewRouteController,
	)
	return &controller.RouteController{}, nil
}

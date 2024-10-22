package biz

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Route struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

type Response struct {
	Code   int     `json:"code"`
	Count  int     `json:"count"`
	Routes []Route `json:"routes"`
}

type RouteBiz struct{}

func NewRouteBiz() *RouteBiz {
	return &RouteBiz{}
}

func (rb *RouteBiz) FetchAndFilterRoutes() (*Response, error) {
	resp, err := http.Get("http://142.171.211.71:6688/all")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	filteredRoutes := make([]Route, 0)
	for _, route := range response.Routes {
		if route.Path != "" {
			filteredRoutes = append(filteredRoutes, route)
		}
	}

	response.Routes = filteredRoutes
	response.Count = len(filteredRoutes)

	return &response, nil
}

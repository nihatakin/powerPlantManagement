package controllers

import (
	"github.com/nihatakin/powerPlantManagement/api/middlewares"
	"net/http"
)

type Route struct {
	Uri          string
	Method       string
	Handler      func(http.ResponseWriter, *http.Request)
	AuthRequired bool
}

func (s *Server) initializeRoutes() {
	for _, route := range Load(s) {
		if route.AuthRequired {
			s.Router.HandleFunc(route.Uri,
				middlewares.SetMiddlewareLogger(
					middlewares.SetMiddlewareJSON(
						middlewares.SetMiddlewareAuthentication(route.Handler))),
			).Methods(route.Method)
		} else {
			s.Router.HandleFunc(route.Uri,
				middlewares.SetMiddlewareLogger(
					middlewares.SetMiddlewareJSON(route.Handler)),
			).Methods(route.Method)
		}
	}
}

func Load(s *Server) []Route {
	routes := HomeRoutes(s)
	routes = append(routes, LoginRoutes(s)...)
	routes = append(routes, UserRoutes(s)...)
	routes = append(routes, PowerPlantRoutes(s)...)
	return routes
}

func HomeRoutes(s *Server) []Route {
	routes := []Route{
		Route{
			Uri:          "/",
			Method:       http.MethodGet,
			Handler:      s.Login,
			AuthRequired: false,
		},
	}
	return routes
}

func LoginRoutes(s *Server) []Route {
	routes := []Route{
		Route{
			Uri:          "/login",
			Method:       http.MethodPost,
			Handler:      s.Login,
			AuthRequired: false,
		},
	}
	return routes
}

func UserRoutes(s *Server) []Route {
	routes := []Route{
		Route{
			Uri:          "/users",
			Method:       http.MethodGet,
			Handler:      s.GetUsers,
			AuthRequired: false,
		},
		Route{
			Uri:          "/users",
			Method:       http.MethodPost,
			Handler:      s.CreateUser,
			AuthRequired: false,
		},
		Route{
			Uri:          "/users/{id}",
			Method:       http.MethodGet,
			Handler:      s.GetUser,
			AuthRequired: false,
		},
		Route{
			Uri:          "/users/{id}",
			Method:       http.MethodPut,
			Handler:      s.UpdateUser,
			AuthRequired: true,
		},
		Route{
			Uri:          "/users/{id}",
			Method:       http.MethodDelete,
			Handler:      s.DeleteUser,
			AuthRequired: true,
		},
	}
	return routes
}

func PowerPlantRoutes(s *Server) []Route {
	routes := []Route{
		Route{
			Uri:          "/powerPlants",
			Method:       http.MethodPost,
			Handler:      s.CreatePowerPlant,
			AuthRequired: false,
		},
		Route{
			Uri:          "/powerPlants",
			Method:       http.MethodGet,
			Handler:      s.GetPowerPlants,
			AuthRequired: false,
		},
		Route{
			Uri:          "/powerPlants/{id}",
			Method:       http.MethodGet,
			Handler:      s.GetPowerPlant,
			AuthRequired: false,
		},
		Route{
			Uri:          "/powerPlants/{id}",
			Method:       http.MethodPut,
			Handler:      s.UpdatePowerPlant,
			AuthRequired: true,
		},
		Route{
			Uri:          "/powerPlants/{id}",
			Method:       http.MethodDelete,
			Handler:      s.DeletePowerPlant,
			AuthRequired: true,
		},
	}
	return routes
}
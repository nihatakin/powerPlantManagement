package controllers

import "github.com/nihatakin/powerPlantManagement/api/middlewares"

func (s *Server) initializeRoutes() {

	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET") // Home Route

	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST") // Login Route

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//PowerPlant routes
	s.Router.HandleFunc("/powerPlants", middlewares.SetMiddlewareJSON(s.CreatePowerPlant)).Methods("POST")
	s.Router.HandleFunc("/powerPlants", middlewares.SetMiddlewareJSON(s.GetPowerPlants)).Methods("GET")
	s.Router.HandleFunc("/powerPlants/{id}", middlewares.SetMiddlewareJSON(s.GetPowerPlant)).Methods("GET")
	s.Router.HandleFunc("/powerPlants/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePowerPlant))).Methods("PUT")
	s.Router.HandleFunc("/powerPlants/{id}", middlewares.SetMiddlewareAuthentication(s.DeletePowerPlant)).Methods("DELETE")

}
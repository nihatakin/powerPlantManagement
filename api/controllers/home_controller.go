package controllers

import (
	"net/http"

	"github.com/nihatakin/powerPlantManagement/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Power Plant Management API")
}
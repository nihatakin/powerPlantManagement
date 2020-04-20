package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nihatakin/powerPlantManagement/api/auth"
	"github.com/nihatakin/powerPlantManagement/api/models"
	"github.com/nihatakin/powerPlantManagement/api/responses"
	"github.com/nihatakin/powerPlantManagement/api/utils/formaterror"
)

func (server *Server) CreatePowerPlant(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	powerPlant := models.PowerPlant{}
	err = json.Unmarshal(body, &powerPlant)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	powerPlant.PrepareForCreate()
	err = powerPlant.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != powerPlant.CreatorUserId {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	powerPlantCreated, err := powerPlant.SavePowerPlant(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, powerPlantCreated.ID))
	responses.JSON(w, http.StatusCreated, powerPlantCreated)
}

func (server *Server) GetPowerPlants(w http.ResponseWriter, r *http.Request) {

	powerPlant := models.PowerPlant{}

	powerPlants, err := powerPlant.FindAllPowerPlants(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, powerPlants)
}

func (server *Server) GetPowerPlant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	powerPlant := models.PowerPlant{}

	powerPlantReceived, err := powerPlant.FindPowerPlantByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, powerPlantReceived)
}

//TODO Cok fazla DB transaction var, iyileştirilmesi lazım
func (server *Server) UpdatePowerPlant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the powerPlant id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the powerPlant exist
	powerPlant := models.PowerPlant{}
	err = server.DB.Debug().Model(models.PowerPlant{}).Where("id = ?", pid).Take(&powerPlant).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("PowerPlant not found"))
		return
	}

	// If a user attempt to update a powerPlant not belonging to him
	/*if uid != powerPlant.LastModifierUserId {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}*/
	// Read the data powerPlanted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	powerPlantUpdate := models.PowerPlant{}
	err = json.Unmarshal(body, &powerPlantUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != powerPlantUpdate.LastModifierUserId {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	powerPlantUpdate.PrepareForUpdate()
	err = powerPlantUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	powerPlantUpdate.ID = powerPlant.ID //this is important to tell the model the powerPlant id to update, the other update field are set above

	powerPlantUpdated, err := powerPlantUpdate.UpdateAPowerPlant(server.DB)

	powerPlantReceived, err := powerPlantUpdate.FindPowerPlantByID(server.DB, powerPlantUpdated.ID)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, powerPlantReceived)
}

func (server *Server) DeletePowerPlant(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid powerPlant id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the powerPlant exist
	powerPlant := models.PowerPlant{}
	err = server.DB.Debug().Model(models.PowerPlant{}).Where("id = ?", pid).Take(&powerPlant).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this powerPlant?
	if uid != powerPlant.CreatorUserId {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = powerPlant.DeleteAPowerPlant(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
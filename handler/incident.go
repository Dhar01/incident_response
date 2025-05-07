package handler

import (
	"net/http"
	"time"

	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/internal/model"

	log "github.com/sirupsen/logrus"
)

func CreateIncident(incident model.IncidentReq, authID uint64) (httpResponse model.HTTPResponse, httpStatusCode int) {
	db := database.GetDB()

	if incident.Title == "" || incident.Status == "" || incident.Severity == "" {
		httpResponse.Message = "Missing required fields"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// check if assignee exists
	if err := db.First(&model.Auth{}, incident.AssignedTo).Error; err != nil {
		log.WithError(err).Error("error code: 2001.1")
		return setErrorMessage("assigned user not found", http.StatusNotFound)
	}

	newIncident := model.Incident{
		Title:       incident.Title,
		Description: incident.Description,
		Status:      incident.Status,
		Severity:    incident.Severity,
		AssignedTo:  incident.AssignedTo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(&newIncident).Error; err != nil {
		log.WithError(err).Error("error code: 2001.2")
		return setErrorMessage(errInternalServer, http.StatusInternalServerError)
	}

	httpResponse.Message = newIncident
	httpStatusCode = http.StatusOK
	return
}

func UpdateIncident(incident model.IncidentUpdate, incidentID uint64) (httpResponse model.HTTPResponse, httpStatusCode int) {
	db := database.GetDB()

	if incident.IncidentID == 0 {
		return setErrorMessage("incident ID is required", http.StatusBadRequest)
	}

	var existing model.Incident

	if err := db.First(&existing, incidentID).Error; err != nil {
		log.WithError(err).Error("error code: 2002.1")
		return setErrorMessage("incident not found", http.StatusNotFound)
	}

	// Update fields
	existing.Title = incident.Title
	existing.Description = incident.Description
	existing.Status = incident.Status
	existing.Severity = incident.Severity
	existing.AssignedTo = incident.AssignedTo
	existing.UpdatedAt = time.Now()

	if err := db.Save(&existing).Error; err != nil {
		log.WithError(err).Error("error code: 2002.2")
		return setErrorMessage(errInternalServer, http.StatusInternalServerError)
	}

	httpResponse.Message = existing
	httpStatusCode = http.StatusOK
	return
}

func GetIncidentByID(id uint64) (httpResponse model.HTTPResponse, httpStatusCode int) {
	db := database.GetDB()

	var incident model.Incident

	if err := db.First(&incident, id).Error; err != nil {
		log.WithError(err).Error("error code: 2003.1")
		return setErrorMessage("incident not found", http.StatusNotFound)
	}

	httpResponse.Message = incident
	httpStatusCode = http.StatusOK
	return
}

func GetAllIncidents() (httpResponse model.HTTPResponse, httpStatusCode int) {
	db := database.GetDB()

	var incidents []model.Incident

	if err := db.Find(&incidents).Error; err != nil {
		log.WithError(err).Error("error code: 2004.1")
		return setErrorMessage(errInternalServer, http.StatusInternalServerError)
	}

	httpResponse.Message = incidents
	httpStatusCode = http.StatusOK
	return
}

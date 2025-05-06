package handler

import (
	"net/http"
	"time"

	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/internal/model"

	log "github.com/sirupsen/logrus"
)

func CreateIncident(incident model.IncidentReq) (httpResponse model.HTTPResponse, httpStatusCode int) {
	db := database.GetDB()

	if incident.Title == "" || incident.Status == "" || incident.Severity == "" {
		httpResponse.Message = "Missing required fields"
		httpStatusCode = http.StatusBadRequest
		return
	}

	// check if assignee exists
	if err := db.First(&model.Auth{}, incident.AssignedTo).Error; err != nil {
		log.WithError(err).Error("error code: 2001.1")
		httpResponse.Message = "Assigned user not found"
		httpStatusCode = http.StatusNotFound
		return
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
		httpResponse.Message = errInternalServer
		httpStatusCode = http.StatusInternalServerError
		return
	}

	httpResponse.Message = newIncident
	httpStatusCode = http.StatusOK
	return
}

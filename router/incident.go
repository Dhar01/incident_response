package router

import (
	"net/http"
	"reflect"

	"github.com/Dhar01/incident_resp/handler"
	"github.com/Dhar01/incident_resp/internal/model"
	"github.com/Dhar01/incident_resp/lib/renderer"
	incident_gen "github.com/Dhar01/incident_resp/router/incidents"
	"github.com/gin-gonic/gin"
)

type incidentAPI struct{}

var _ incident_gen.ServerInterface = (*incidentAPI)(nil)

func newIncidentAPI() *incidentAPI {
	return &incidentAPI{}
}

func (api *incidentAPI) FetchIncidents(c *gin.Context) {

}

func (api *incidentAPI) CreateNewIncident(c *gin.Context) {
	// ! user verifications will be processed on middleware
	var req model.IncidentReq

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		renderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.CreateIncident(req)

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		renderer.Render(c, resp, statusCode)
		return
	}

	renderer.Render(c, resp.Message, statusCode)
}

func (api *incidentAPI) FetchIncidentByID(c *gin.Context, id uint64) {

}

func (api *incidentAPI) UpdateIncident(c *gin.Context, id uint64) {

}

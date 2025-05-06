package router

import (
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

}

func (api *incidentAPI) FetchIncidentByID(c *gin.Context, id uint64) {

}

func (api *incidentAPI) UpdateIncident(c *gin.Context, id uint64) {

}

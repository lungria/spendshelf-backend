package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/src/report"
	"go.uber.org/zap"
)

type ReportsHandler struct {
	generator report.Generator
	logger    *zap.SugaredLogger
}

type GetQuery struct {
	Start time.Time `form:"start" time_format:"unix"`
	End   time.Time `form:"end" time_format:"unix"`
}

type GetResponse struct {
	Report []report.Element `json:"report"`
}

func NewReportsHandler(generator report.Generator, logger *zap.SugaredLogger) *ReportsHandler {
	return &ReportsHandler{generator: generator, logger: logger}
}

// HandleGet return all existing categories
func (handler *ReportsHandler) HandleGet(c *gin.Context) {
	var query GetQuery
	err := c.BindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseFromError(err, "Unable to parse query"))
		return
	}
	reportResponse, err := handler.generator.GetReport(c, query.Start, query.End)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseFromError(err, "Unable to generate report"))
		return
	}
	c.JSON(http.StatusOK, GetResponse{reportResponse})
}

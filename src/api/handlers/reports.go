package handlers

import (
	"errors"
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

type getQuery struct {
	From time.Time `form:"from" time_format:"unix"`
	To   time.Time `form:"to" time_format:"unix"`
}

type getResponse struct {
	Report []report.Element `json:"report"`
}

func NewReportsHandler(generator report.Generator, logger *zap.SugaredLogger) *ReportsHandler {
	return &ReportsHandler{generator: generator, logger: logger}
}

// HandleGet return all existing categories
func (handler *ReportsHandler) HandleGet(c *gin.Context) {
	var query getQuery
	err := c.BindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFromError(err, "Unable to parse query"))
		return
	}
	if query.To.Before(query.From) || query.To.Equal(query.From) {
		c.JSON(http.StatusBadRequest, responseFromError(errors.New("wrong_date_limits"), "To must be after From"))
	}
	reportResponse, err := handler.generator.GetReport(c, query.From, query.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseFromError(err, "Unable to generate report"))
		return
	}
	c.JSON(http.StatusOK, getResponse{reportResponse})
}

package services

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ginHandler struct {
	Handler
}

// GinRoutes sets routes for http.ListenAndServe.
func GinRoutes(h Handler) *gin.Engine {
	hgin := ginHandler{Handler: h}
	r := gin.Default()
	r.GET("/api/visits", hgin.getVisitsHandler)
	r.POST("/api/visits", hgin.postVisitHandler)
	return r
}

func (h ginHandler) getVisitsHandler(c *gin.Context) {
	filter := make(map[string]string)
	for f, val := range c.Request.URL.Query() {
		filter[f] = val[0]
	}
	visits, err := h.GetVisits(filter)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "unexpected error occured"},
		)
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Writer).Encode(visits); err != nil {
		http.Error(c.Writer, "unexpected error occured", http.StatusInternalServerError)
		return
	}
}

func (h ginHandler) postVisitHandler(c *gin.Context) {
	ip := strings.Split(c.Request.RemoteAddr, ":")[0]
	if err := h.ProduceVisit(ip); err != nil {
		http.Error(c.Writer, "unexpected error occured", http.StatusInternalServerError)
		return
	}
}

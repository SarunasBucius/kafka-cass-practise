package services

import (
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
	r.GET("/api/visits/:ip", hgin.getVisitsByIPHandler)
	return r
}

func (h ginHandler) getVisitsHandler(c *gin.Context) {
	filter := make(map[string]string)
	for f, val := range c.Request.URL.Query() {
		filter[f] = val[0]
	}
	visits, err := h.GetVisits(filter)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(200, visits)
}

func (h ginHandler) postVisitHandler(c *gin.Context) {
	ip := strings.Split(c.Request.RemoteAddr, ":")[0]
	if err := h.ProduceVisit(ip); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h ginHandler) getVisitsByIPHandler(c *gin.Context) {
	ip := c.Param("ip")
	filter := make(map[string]string)
	for f, val := range c.Request.URL.Query() {
		filter[f] = val[0]
	}
	visits, err := h.GetVisitsByIP(ip, filter)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(200, visits)
}

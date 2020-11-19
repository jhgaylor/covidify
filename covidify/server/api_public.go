/*
 * Covidify
 *
 * Simple API collecting guest data.
 *
 * API version: 1.0.0
 * Contact: you@your-company.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package covidify

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatz/covidify/covidify/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddVisit - adds an Visit entry
func (s *Server) AddVisit(c *gin.Context) {
	var visit models.Visit
	if err := c.ShouldBindJSON(&visit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := uuid.NewRandom()
	if err == nil {
		visit.Id = u.String()
	}

	if visit.CheckIn.Unix() <= 0 {
		visit.CheckIn = time.Now()
	}

	if err := visit.Valid(); err != nil {
		s.config.Logger.Errorf("Visit Data invalid: %v - %s", visit, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Visit Data not valid: %s", err.Error())})
		return
	}

	v, err := s.db.CreateVisit(visit)
	if err != nil {
		s.config.Logger.Errorf("Error storing visit %v - %s", visit, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not store visit"})
		return
	}

	//StatsD
	s.statsDIncrement("covidify.visit")
	s.statsDIncrement(fmt.Sprintf("covidify.visit.table.%s", v.TableNumber))
	s.statsDIncrementByValue("covidify.visitors", len(v.Visitors))

	//Prometheus

	c.JSON(http.StatusCreated, v)
}

// CheckVisit - Visit status check
func (s *Server) CheckVisit(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

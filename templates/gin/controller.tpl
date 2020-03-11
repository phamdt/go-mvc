package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/boil"
	"go.uber.org/zap"
)

{{#Operations}}
// {{Name}}Controller exposes the methods for interacting with the
// RESTful {{Name}} resource
type {{Name}}Controller struct {
	db  *sqlx.DB
	log *zap.Logger
}

{{#GET}}
// Index returns a list of {{Name}} records
func (ctrl *{{Name}}Controller) Index(c *gin.Context) {
	q := c.Request.URL.RawQuery
	qms := GetQueryModFromQuery(q)
	results, err := models.{{Name}}(qms...).All(ctrl.db)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, results)
}{{/GET}}
{{#POST}}
// Create saves a new {{Name}} record into the database
func (ctrl *{{Name}}Controller) Create(c *gin.Context) {
	m := models.{{Name}}{}
	if err := c.ShouldBindJSON(m); err != nil {
		ctrl.log.Error("invalid {{Name}} creation request",
			zap.Error(err),
		)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := m.Insert(ctrl.db, boil.Infer())
	if err != nil {
		ctrl.log.Error("error creating {{Name}}",
			zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusCreated, gin.H{})
}{{/POST}}
{{#GET}}
// Show retrieves a new {{Name}} record from the database
func (ctrl *{{Name}}Controller) Show(c *gin.Context) {
	m := models.{{Name}}{}
	if err := c.ShouldBindUri(&m); err != nil {
		ctrl.log.Error("invalid {{Name}} retrieval request",
			zap.Error(err),
		)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	result, err := models.Find{{Name}}(id)
	if err != nil {
		ctrl.log.Error("error retrieving {{Name}}",
			zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, result)
}{{/GET}}

{{#PUT}}
// Update updates a new {{Name}} record in the database
func (ctrl *{{Name}}Controller) Update(c *gin.Context) {
	m := models.{{Name}}{}
	if err := c.ShouldBindUri(&m); err != nil {
		ctrl.log.Error("invalid {{Name}} update request",
			zap.Error(err),
		)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBindJSON(&m); err != nil {
		ctrl.log.Error("invalid {{Name}} update request",
			zap.Error(err),
		)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := m.Update(ctrl.db, boil.Infer())
	if err != nil {
		ctrl.log.Error("error updating {{Name}}",
			zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, gin.H{})
}{{/PUT}}
{{#DELETE}}
// Delete deletes a new {{Name}} record into the database
func (ctrl *{{Name}}Controller) Delete(c *gin.Context) {
	m := models.{{Name}}{}
	if err := c.ShouldBindUri(&m); err != nil {
		ctrl.log.Error("invalid {{Name}} deletion request",
			zap.Error(err),
		)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err := m.Delete(ctrl.db)
	if err != nil {
		ctrl.log.Error("error deleting {{Name}}",
			zap.Error(err))
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.JSON(http.StatusOK, gin.H{})
}{{/DELETE}}
{{/Operations}}
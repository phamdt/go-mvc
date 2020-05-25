package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func GetRouter(log *zap.SugaredLogger, db *sqlx.DB) *gin.Engine {
	r := gin.New()

{{#Controllers}}
	{{Name}}Ctrl := {{Name}}Controller{db: db, log: log}
{{#Operations}}
	r.{{Method}}("{{Path}}", {{Name}}Ctrl.{{Handler}})
{{/Operations}}
{{/Controllers}}
	return r
}

package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"{{ModuleName}}/models"

	"github.com/jmoiron/sqlx"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

{{#each Actions}}
{{{ whichActionTest Handler }}}
{{/each}}

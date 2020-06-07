package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

{{#each Actions}}
{{{ whichActionTest Handler }}}
{{/each}}

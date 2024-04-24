package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kareemmahlees/meta-x/models"
	"github.com/kareemmahlees/meta-x/utils"
	"github.com/stretchr/testify/assert"
)

func TestHttpError(t *testing.T) {
	rr := httptest.NewRecorder()
	httpError(rr, http.StatusOK, "something")

	assert.Equal(t, rr.Code, http.StatusOK)

	body := utils.DecodeBody[models.ErrResp](rr.Result().Body)
	assert.Equal(t, body.Message, "something")
}

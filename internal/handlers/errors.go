package handlers

import (
	"net/http"

	"github.com/kareemmahlees/meta-x/models"
)

func httpError(w http.ResponseWriter, code int, errMsg any) {
	w.WriteHeader(code)

	writeJson(w, models.ErrResp{
		Code:    code,
		Message: errMsg,
	})
}

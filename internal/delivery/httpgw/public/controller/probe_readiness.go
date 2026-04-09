package controller

import (
	"io"
	"net/http"
)

func (ctrl *Controller) ProbeReadiness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

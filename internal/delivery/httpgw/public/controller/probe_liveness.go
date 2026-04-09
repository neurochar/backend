package controller

import (
	"io"
	"net/http"
)

func (ctrl *Controller) ProbeLiveness(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}

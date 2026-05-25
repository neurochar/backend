package private_http

import (
	"encoding/json"
	"net/http"
)

type metricsResponse struct {
	Value int `json:"value"`
}

func (s *Server) MetricsShouldGPUNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metricsResponse{Value: 0})
}

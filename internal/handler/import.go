package handler

import (
	"encoding/json"
	"net/http"
)

type ImportRequest struct {
	Numbers []string `json:"numbers"`
	Source  string   `json:"source"`
}

type ImportResponse struct {
	Status string `json:"status"`
}

func (h *Handler) Import(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Source == "" {
		http.Error(w, "source is required", http.StatusBadRequest)
		return
	}

	if len(req.Numbers) == 0 {
		http.Error(w, "numbers array is empty", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ImportResponse{Status: "ok"})
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IlushaSPB/test-phone-number-service/internal/db"
	"github.com/IlushaSPB/test-phone-number-service/internal/service"
)

type ImportRequest struct {
	Numbers []string `json:"numbers"`
	Source  string   `json:"source"`
}

type ImportResponse struct {
	Accepted int      `json:"accepted"`
	Skipped  int      `json:"skipped"`
	Errors   int      `json:"errors"`
	Details  []string `json:"details,omitempty"`
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

	ctx := context.Background()
	resp := ImportResponse{}

	seen := make(map[string]bool)

	for _, raw := range req.Numbers {
		info, err := service.NormalizeAndEnrich(raw)
		if err != nil {
			resp.Errors++
			resp.Details = append(resp.Details, raw+": "+err.Error())
			continue
		}

		if seen[info.E164] {
			resp.Skipped++
			continue
		}
		seen[info.E164] = true

		rows, err := h.queries.InsertPhoneNumber(ctx, db.InsertPhoneNumberParams{
			PhoneNumber: info.E164,
			Source:      req.Source,
			Country:     info.Country,
			Region:      info.Region,
			Provider:    info.Provider,
		})

		if err != nil {
			resp.Errors++
			resp.Details = append(resp.Details, info.E164+": db error")
		} else if rows == 0 {
			resp.Skipped++
		} else {
			resp.Accepted++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/IlushaSPB/test-phone-number-service/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type SearchResponse struct {
	Total   int64         `json:"total"`
	Limit   int           `json:"limit"`
	Offset  int           `json:"offset"`
	Numbers []PhoneNumber `json:"numbers"`
}

type PhoneNumber struct {
	ID          int64  `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Source      string `json:"source"`
	Country     string `json:"country"`
	Region      string `json:"region,omitempty"`
	Provider    string `json:"provider,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	number := query.Get("number")
	country := query.Get("country")
	region := query.Get("region")
	provider := query.Get("provider")

	limit := 20
	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
			if limit > 100 {
				limit = 100
			}
		}
	}

	offset := 0
	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	ctx := context.Background()

	countParams := db.CountPhoneNumbersParams{
		Number:   toNullText(number),
		Country:  toNullText(country),
		Region:   toNullText(region),
		Provider: toNullText(provider),
	}

	total, err := h.queries.CountPhoneNumbers(ctx, countParams)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	searchParams := db.SearchPhoneNumbersParams{
		Number:   toNullText(number),
		Country:  toNullText(country),
		Region:   toNullText(region),
		Provider: toNullText(provider),
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	rows, err := h.queries.SearchPhoneNumbers(ctx, searchParams)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	numbers := make([]PhoneNumber, 0, len(rows))
	for _, row := range rows {
		createdAt := ""
		if row.CreatedAt.Valid {
			createdAt = row.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
		}
		numbers = append(numbers, PhoneNumber{
			ID:          row.ID,
			PhoneNumber: row.PhoneNumber,
			Source:      row.Source,
			Country:     row.Country,
			Region:      nullTextToString(row.Region),
			Provider:    nullTextToString(row.Provider),
			CreatedAt:   createdAt,
		})
	}

	resp := SearchResponse{
		Total:   total,
		Limit:   limit,
		Offset:  offset,
		Numbers: numbers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func toNullText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}

func nullTextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

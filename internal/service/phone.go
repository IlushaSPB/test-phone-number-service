package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

var digitRegex = regexp.MustCompile(`[^\d+]`)

type PhoneInfo struct {
	E164     string
	Country  string
	Region   pgtype.Text
	Provider pgtype.Text
}

func NormalizeAndEnrich(raw string) (*PhoneInfo, error) {
	normalized, err := normalizeE164(raw)
	if err != nil {
		return nil, err
	}

	country := detectCountry(normalized)
	if country == "" {
		return nil, fmt.Errorf("unknown country code")
	}

	var region, provider pgtype.Text
	if strings.HasPrefix(normalized, "+7") && len(normalized) == 12 {
		r, p := detectRussianRegionProvider(normalized)
		if r != "" {
			region = pgtype.Text{String: r, Valid: true}
		}
		if p != "" {
			provider = pgtype.Text{String: p, Valid: true}
		}
	}

	return &PhoneInfo{
		E164:     normalized,
		Country:  country,
		Region:   region,
		Provider: provider,
	}, nil
}

func normalizeE164(raw string) (string, error) {
	cleaned := digitRegex.ReplaceAllString(raw, "")

	if cleaned == "" {
		return "", fmt.Errorf("no digits found")
	}

	if strings.HasPrefix(cleaned, "+8") && len(cleaned) == 12 {
		cleaned = "+7" + cleaned[2:]
	} else if strings.HasPrefix(cleaned, "8") && len(cleaned) == 11 {
		cleaned = "+7" + cleaned[1:]
	} else if strings.HasPrefix(cleaned, "7") && len(cleaned) == 11 {
		cleaned = "+" + cleaned
	} else if !strings.HasPrefix(cleaned, "+") {
		if len(cleaned) == 10 {
			cleaned = "+7" + cleaned
		} else {
			cleaned = "+" + cleaned
		}
	}

	if len(cleaned) < 8 || len(cleaned) > 16 {
		return "", fmt.Errorf("invalid length: %d", len(cleaned))
	}

	return cleaned, nil
}

func detectCountry(e164 string) string {
	countryMap := map[string]string{
		"+7":   "Russia",
		"+1":   "USA",
		"+44":  "UK",
		"+49":  "Germany",
		"+33":  "France",
		"+86":  "China",
		"+81":  "Japan",
		"+82":  "South Korea",
		"+91":  "India",
		"+55":  "Brazil",
		"+375": "Belarus",
		"+998": "Uzbekistan",
	}

	for prefix, country := range countryMap {
		if strings.HasPrefix(e164, prefix) {
			return country
		}
	}

	return ""
}

func detectRussianRegionProvider(e164 string) (region, provider string) {
	if len(e164) != 12 || !strings.HasPrefix(e164, "+7") {
		return "", ""
	}

	code := e164[2:5]

	regionProviderMap := map[string]struct {
		Region   string
		Provider string
	}{
		"916": {"Москва", "МТС"},
		"917": {"Москва", "МТС"},
		"915": {"Москва", "МТС"},
		"910": {"Воронеж", "МТС"},
		"919": {"Москва", "МТС"},
		"913": {"Новосибирск", "МТС"},
		"925": {"Москва", "МегаФон"},
		"926": {"Москва", "МегаФон"},
		"920": {"Нижний Новгород", "МегаФон"},
		"922": {"Тюмень", "МегаФон"},
		"927": {"Самара", "МегаФон"},
		"928": {"Дагестан", "МегаФон"},
		"903": {"Москва", "Билайн"},
		"905": {"", "Билайн"},
		"906": {"Ставрополь", "Билайн"},
		"909": {"Самара", "Билайн"},
	}

	if info, ok := regionProviderMap[code]; ok {
		return info.Region, info.Provider
	}

	return "", ""
}

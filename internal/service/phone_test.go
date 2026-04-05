package service

import (
	"testing"
)

func TestNormalizeE164(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "10 цифр без префикса",
			input: "9161234567",
			want:  "+79161234567",
		},
		{
			name:  "начинается с 8",
			input: "89161234567",
			want:  "+79161234567",
		},
		{
			name:  "с пробелами и форматированием",
			input: "+7 (916) 123-45-67",
			want:  "+79161234567",
		},
		{
			name:  "уже в формате E164",
			input: "+79161234567",
			want:  "+79161234567",
		},
		{
			name:  "номер США",
			input: "+12025551234",
			want:  "+12025551234",
		},
		{
			name:    "слишком короткий",
			input:   "123",
			wantErr: true,
		},
		{
			name:    "нет цифр",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeE164(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("normalizeE164() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("normalizeE164() unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeE164() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectCountry(t *testing.T) {
	tests := []struct {
		e164 string
		want string
	}{
		{"+79161234567", "Russia"},
		{"+12025551234", "USA"},
		{"+441234567890", "UK"},
		{"+491234567890", "Germany"},
		{"+33123456789", "France"},
		{"+861234567890", "China"},
		{"+999999999", ""},
	}

	for _, tt := range tests {
		t.Run(tt.e164, func(t *testing.T) {
			got := detectCountry(tt.e164)
			if got != tt.want {
				t.Errorf("detectCountry(%s) = %v, want %v", tt.e164, got, tt.want)
			}
		})
	}
}

func TestDetectRussianRegionProvider(t *testing.T) {
	tests := []struct {
		e164         string
		wantRegion   string
		wantProvider string
	}{
		{"+79161234567", "Москва", "МТС"},
		{"+79251234567", "Москва", "МегаФон"},
		{"+79031234567", "Москва", "Билайн"},
		{"+79221234567", "Тюмень", "МегаФон"},
		{"+79101234567", "Воронеж", "МТС"},
		{"+79051234567", "", "Билайн"},
		{"+79991234567", "", ""},
		{"+12025551234", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.e164, func(t *testing.T) {
			gotRegion, gotProvider := detectRussianRegionProvider(tt.e164)
			if gotRegion != tt.wantRegion {
				t.Errorf("detectRussianRegionProvider(%s) region = %v, want %v", tt.e164, gotRegion, tt.wantRegion)
			}
			if gotProvider != tt.wantProvider {
				t.Errorf("detectRussianRegionProvider(%s) provider = %v, want %v", tt.e164, gotProvider, tt.wantProvider)
			}
		})
	}
}

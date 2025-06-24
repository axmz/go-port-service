package port

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		portName string
		city     string
		country  string
		wantErr  bool
	}{
		{"all fields valid", "id", "name", "city", "country", false},
		{"missing id", "", "name", "city", "country", true},
		{"missing name", "id", "", "city", "country", true},
		{"missing city", "id", "name", "", "country", true},
		{"missing country", "id", "name", "city", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.id, tt.portName, tt.city, tt.country)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package port

import "time"

type Port struct {
	ID          string
	Name        string
	Code        string
	City        string
	Country     string
	Alias       []string
	Regions     []string
	Coordinates []float64
	Province    string
	Timezone    string
	Unlocs      []string

	// TODO: use these fields.
	CreatedAt time.Time
	UpdatedAt time.Time
}

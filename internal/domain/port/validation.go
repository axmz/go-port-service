package port

import "fmt"

func validate(id, name, city, country string) error {
	// TODO: add validation library?
	fields := map[string]string{
		"port id":      id,
		"port name":    name,
		"port city":    city,
		"port country": country,
	}

	for field, value := range fields {
		if value == "" {
			return fmt.Errorf("%w: %s", ErrRequired, field)
		}
	}

	return nil
}

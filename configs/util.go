package configs

import (
	"fmt"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}
	d.Duration = duration
	return nil
}

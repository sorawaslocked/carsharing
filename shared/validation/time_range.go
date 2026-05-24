package validation

import "time"

type TimeRange struct {
	From *time.Time
	To   *time.Time
}

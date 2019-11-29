package monitor

import (
	"fmt"
	"github.com/cultureamp/glamplify/field"
)

// Entries contains key-value pairs to record along with the event
type Fields field.Fields

// Validate checks that Entries are valid before processing
func (fields Fields) Validate() (bool, error) {
	// https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits

	for k, v := range fields {

		switch s := v.(type) {
		case nil:
			return false, fmt.Errorf("key '%v' cannot have 'nil' value", k)
		case string:
			if len(s) > 254 {
				return false, fmt.Errorf("key '%v' too long, must be less than 255 characters", k)
			}
		case float32, float64, int32, int64, int:
			continue
		default:
			return false, fmt.Errorf("key '%v' must be string, float or int data type", k)
		}
	}

	return true, nil
}

package log

import (
	"encoding/json"
	"fmt"
)

// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

func (fields Fields) Merge(other ...Fields) Fields {
	merged := Fields{}

	for k, v := range fields {
		merged[k] = v
	}

	for _, f := range other {
		for k, v := range f {
			merged[k] = v
		}
	}

	return merged
}

func (fields Fields) Serialize() string {

	bytes, err := json.Marshal(fields)
	if err != nil {
		// REVISIT - panic?
	}

	return string(bytes)
}

// Validate checks that Entries are valid before processing
func (fields Fields) ValidateNewRelic() (bool, error) {
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

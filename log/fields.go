package log

import (
	"encoding/json"
	"fmt"
	"github.com/cultureamp/glamplify/helper"
	systemLog "log"
	"reflect"
	"time"
)


// Fields type, used to pass to Debug, Print and Error.
type Fields map[string]interface{}

func NewDurationFields(duration time.Duration) Fields {
	return Fields{
		TimeTaken: DurationAsISO8601(duration),
		TimeTakenMS: duration.Milliseconds(),
	}
}

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

func (fields Fields) ToSnakeCase() Fields {
	snaked := Fields{}

	for k, v := range fields {

		switch f := v.(type) {
		case Fields:
			v = f.ToSnakeCase()
		}

		sc := helper.ToSnakeCase(k)
		snaked[sc] = v
	}

	return snaked
}

func (fields Fields) ToJson() string {

	filtered := fields.filterNonSerializableValues()
	bytes, err := json.Marshal(filtered)
	if err != nil {
		systemLog.Printf("failed to serialize log fields to json string. err: %s", err.Error())
		// REVISIT - panic?
	}

	return string(bytes)
}

// ValidateNewRelic checks that Entries are valid according to NewRelic requirements before processing
// https://docs.newrelic.com/docs/insights/insights-data-sources/custom-data/insights-custom-data-requirements-limits
func (fields Fields) ValidateNewRelic() (bool, error) {

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

func (fields Fields) filterNonSerializableValues() Fields {
	filtered := Fields{}

	for k, v := range fields {
		vt := reflect.TypeOf(v).Kind()
		switch vt {
		case reflect.Func, reflect.Chan : // add other types we don't want to log here
			continue
		default:
			filtered[k] = v
		}
	}
	return filtered
}


package common

import (
	"reflect"

	"github.com/pkg/errors"
)

var modelIdPrefixes = map[string]string{
	"User":         "user",
	"Platform":     "platform",
	"Network":      "network",
	"Asset":        "asset",
	"Device":       "device",
	"Contact":      "contact",
	"Location":     "location",
	"Instrument":   "instrument",
	"TxLeg":        "txleg",
	"Transaction":  "tx",
	"AuthStrategy": "auth",
	"ApiKey":       "apikey",
	"Contract":     "contract",
}

var relationalIdPrefixes = map[string]string{
	"UserId":             "user",
	"PlatformId":         "platform",
	"ContactId":          "contact",
	"DeviceId":           "device",
	"InstrumentId":       "instrument",
	"NetworkId":          "network",
	"OriginTxLegId":      "txleg",
	"DestinationTxLegId": "txleg",
	"ReceiptTxId":        "txleg",
	"ResponseTxId":       "txleg",
	"AssetId":            "asset",
}

func SanitizeIdInput(model interface{}, inline ...*string) error {
	stype := reflect.ValueOf(model).Elem()

	// Modify struct passed in by value
	if len(inline) == 0 {
		// Model ID
		field := stype.FieldByName("Id")
		if !field.IsValid() {
			// return errors.New("model does not contain an id")
			// model may be relational or inline, continue
		} else {
			prefix, ok := modelIdPrefixes[stype.Type().Name()]
			if !ok {
				return StringError(errors.New("model unknown"))
			}
			if field.String()[:len(prefix)+1] != prefix+"_" {
				return StringError(errors.New("input missing prefix " + prefix + "_"))
			}
			field.SetString(field.String()[len(prefix)+1:])
		}

		// Relational IDs
		// Iterating through const memory is faster than Inline logic below despite ordering of fields
		for fieldName, prefix := range relationalIdPrefixes {
			field := stype.FieldByName(fieldName)
			if !field.IsValid() {
				continue
			}
			if field.String()[:len(prefix)+1] != prefix+"_" {
				return StringError(errors.New("input missing prefix " + prefix + "_"))
			}
			field.SetString(field.String()[len(prefix)+1:])
		}
		return nil
	}

	// Modify inline string pointers based on struct values
	if len(inline) > 0 {
		if len(inline) != stype.NumField() {
			return StringError(errors.New("invalid inline length"))
		}
		for i := 0; i < stype.NumField(); i++ {
			field := stype.Field(i)
			prefix, ok := relationalIdPrefixes[stype.Type().Field(i).Name]
			if !ok {
				return StringError(errors.New("unknown field type " + field.Type().Name()))
			}
			if field.String()[:len(prefix)+1] != prefix+"_" {
				return StringError(errors.New("input missing prefix " + prefix + "_"))
			}
			*inline[i] = field.String()[len(prefix)+1:]
		}
	}

	return nil
}

func SanitizeIdOutput(model interface{}, inline ...*string) error {
	stype := reflect.ValueOf(model).Elem()

	if len(inline) == 0 {
		// Model ID
		field := stype.FieldByName("Id")
		if !field.IsValid() {
			// return errors.New("model does not contain an id")
			// Model may be relational, continue
		} else {
			prefix, ok := modelIdPrefixes[stype.Type().Name()]
			if !ok {
				return StringError(errors.New("model unknown"))
			}
			field.SetString(prefix + "_" + field.String())
		}

		// Relational IDs
		for fieldName, prefix := range relationalIdPrefixes {
			field := stype.FieldByName(fieldName)
			if !field.IsValid() {
				continue
			}
			field.SetString(prefix + "_" + field.String())
		}

		return nil
	}

	// Modify inline string pointers based on struct values
	if len(inline) > 0 {
		if len(inline) != stype.NumField() {
			return StringError(errors.New("invalid inline length"))
		}
		for i := 0; i < stype.NumField(); i++ {
			field := stype.Field(i)
			prefix, ok := relationalIdPrefixes[stype.Type().Field(i).Name]
			if !ok {
				return StringError(errors.New("unknown field type " + field.Type().Name()))
			}
			*inline[i] = prefix + "_" + field.String()
		}
	}

	return nil
}

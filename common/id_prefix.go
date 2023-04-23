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

func SanitizeIdInput(model interface{}) error {
	// Model ID
	stype := reflect.ValueOf(model).Elem()
	field := stype.FieldByName("Id")
	if !field.IsValid() {
		// return errors.New("model does not contain an id")
		// model may be relational, continue
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

func SanitizeIdOutput(model interface{}) error {
	// Model ID
	stype := reflect.ValueOf(model).Elem()
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

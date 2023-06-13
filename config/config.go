package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadEnv(collection any, params ...string) error {
	path := ".env"
	if len(params) >= 1 && params[0] != "" {
		path = params[0]
	}
	godotenv.Load(path)
	missing := []string{}
	stype := reflect.ValueOf(collection).Elem()
	for i := 0; i < stype.NumField(); i++ {
		field := stype.Field(i)
		key := stype.Type().Field(i).Name
		value := os.Getenv(key)
		required := stype.Type().Field(i).Tag.Get("required") == "true"
		optional := stype.Type().Field(i).Tag.Get("required") == "false"
		if required && value == "" {
			missing = append(missing, key)
		}
		if optional && value == "" {
			// lets not panic, but warn
			log.Warn().Str("env var", key).Msg("Optional environment variable not set")
		}
		field.SetString(value)
	}
	if len(missing) > 0 {
		panic("Missing environment variable: " + strings.Join(missing, ", "))
	}
	return nil
}

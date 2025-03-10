package helper

import (
	"log"
	"os"
	"strings"
)

func GetMedia(name string) []string {
	return strings.Split(GetEnvVar(name), ",")
}

func GetEnvVar(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Environment variable '%s' was not defined", name)
	}
	return value
}

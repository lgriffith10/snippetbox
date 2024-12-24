package env

import (
	"log"
	"os"
)

var dev = map[string]string{
	"GO_PORT": ":8080",
}

var prod = map[string]string{
	"GO_PORT": ":8081",
}

func SetEnvVariables(mode string) {
	var keys map[string]string

	switch mode {
	case "dev":
		keys = dev
	case "prod":
		keys = prod
	default:
		log.Fatalf("Unsupported environment mode: %s", mode)
	}

	for key, value := range keys {
		err := os.Setenv(key, value)

		if err != nil {
			log.Fatal(err)
		}
	}
}
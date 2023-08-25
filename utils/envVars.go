package utils

import (
	"github.com/charmbracelet/log"

	"os"
)

func GetEnvVar(key string) string {
	envVar, exists := os.LookupEnv(key)
	if !exists {
		log.Errorf("%s Env Var is missing",key)
		return ""
	}
		return envVar
}
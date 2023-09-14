package utils

import (
	"fmt"
	"os"
)

func GetEnvVar(key string, allowUnset bool) (string, error) {
	envVar, exists := os.LookupEnv(key)

	if !exists && allowUnset {
		return "", nil
	} else if !exists && !allowUnset {
		return "", fmt.Errorf("%s Env Var is missing", key)
	}
	return envVar, nil
}

package utils

import (
	"fmt"
	"os"
)

func GetEnvVar(key string) (string,error) {
	envVar, exists := os.LookupEnv(key)
	if !exists {
		return "",fmt.Errorf("%s Env Var is missing",key)
	}
		return envVar,nil
}
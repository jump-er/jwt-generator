package config

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func LoggerInit() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		DisableColors: false,
		FullTimestamp: true,
	})
}

func GetVaultEnv(path string, envs []string) (string, error) {
	precursiveVaultEnvRaw := strings.Split(path, string(os.PathSeparator))

	if len(precursiveVaultEnvRaw) >= 1 {
		precursiveVaultEnv := precursiveVaultEnvRaw[1]
		for _, env := range envs {
			if precursiveVaultEnv == env {
				return precursiveVaultEnv, nil
			}
		}
	}

	return "", fmt.Errorf("Vault env not get. Unknown vault env.")
}

func GetPrivateKey(vaultEnv, mainPath, privateKeyJwtGenerator, privateKeyDefault string) (string, error) {
	privateKeyJwtGeneratorFullPath := mainPath + "/vault/" + vaultEnv + privateKeyJwtGenerator
	privateKeyDefaultFullPath := mainPath + "/vault/" + vaultEnv + privateKeyDefault

	if _, err := os.Stat(privateKeyJwtGeneratorFullPath); err == nil {
		return privateKeyJwtGeneratorFullPath, nil
	}

	if _, err := os.Stat(privateKeyDefaultFullPath); err == nil {
		return privateKeyDefaultFullPath, nil
	}

	return "", fmt.Errorf("Private key not found. Unknown vault env. Check the path %s and %s", privateKeyJwtGenerator, privateKeyDefault)
}

package main

import (
	"fmt"
	"jwt-generator/cmd"
	"jwt-generator/config"
	"jwt-generator/schema"
	"os"
	"path/filepath"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	PayloadRootDir string = "./data"
	MainPath       string = "MAIN"
	CutSymbolCount int    = 120
)

var Envs []string = []string{"test", "stage", "prod"}
var PayloadFiles []string
var PrivateKeyJwtGenerator string = "/path/to/key.pem"
var PrivateKeyDefault string = "/path/to/default/key.pem"

func main() {
	config.LoggerInit()

	log.Info("Start jwt-generator.")
	cutLine()

	err := filepath.Walk(PayloadRootDir, getPayloads)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, file := range PayloadFiles {
		wg.Add(1)
		go run(&wg, file)
	}
	wg.Wait()

	cutLine()
	log.Info("Complete jwt-generator.")
}

func getPayloads(path string, info os.FileInfo, err error) error {
	if filepath.Ext(path) == ".json" {
		PayloadFiles = append(PayloadFiles, path)
	}
	return nil
}

func run(wg *sync.WaitGroup, file string) {
	defer wg.Done()

	err := schema.ValidatePayloadScheme(file)
	if err != nil {
		log.Errorf("%s - %s", file, err)
	}

	vaultEnv, err := config.GetVaultEnv(file, Envs)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := config.GetPrivateKey(vaultEnv, os.Getenv(MainPath), PrivateKeyJwtGenerator, PrivateKeyDefault)
	if err != nil {
		log.Fatal(err)
	}

	_, err = cmd.CreateToken(file, privateKey)
	if err != nil {
		log.Error(err)
	}
}

func cutLine() {
	cut := strings.Repeat("-", CutSymbolCount)
	fmt.Println(cut)
}

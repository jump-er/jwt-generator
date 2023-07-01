package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

type PayloadSchema struct {
	Data        map[string]interface{} `json:"data" validate:"required"`
	Gates       []string               `json:"gates" validate:"unique"`
	Permissions map[string]interface{} `json:"permissions" validate:"required"`
}

func ValidatePayloadScheme(payloadPath string) error {
	payloadData, err := ioutil.ReadFile(payloadPath)
	if err != nil {
		return fmt.Errorf("File %s read problem.", payloadPath)
	}

	var payloadSchema PayloadSchema
	err = json.Unmarshal(payloadData, &payloadSchema)
	if err != nil {
		return err
	}

	v := validator.New()
	log.Infof("Validating %s...", payloadPath)

	err = v.Struct(payloadSchema)

	return err
}

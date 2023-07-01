package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func AddServiceFields(srcPayload map[string]interface{}, payloadPath string) {
	nowTime := time.Now().Unix()

	srcPayload["iat"] = nowTime
	srcPayload["ser"] = GetMD5Hash(payloadPath)
	srcPayload["jti"] = GetMD5Hash(strconv.FormatInt(nowTime, 10))
}

func CreateToken(payloadJsonPath, privateKeyPath string) (string, error) {
	privateKeyRaw, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("Error reading private key file: %v\n", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyRaw)
	if err != nil {
		return "", fmt.Errorf("Error parsing private key file: %v\n", err)
	}

	log.Infof("Preparation data for token %s...", payloadJsonPath)
	jsonStr, err := ioutil.ReadFile(payloadJsonPath)
	if err != nil {
		return "", err
	}

	var payloadPrepare map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &payloadPrepare)
	if err != nil {
		return "", err
	}

	AddServiceFields(payloadPrepare, payloadJsonPath)

	payload := jwt.MapClaims(payloadPrepare)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	var tokenStatus string
	if _, err := os.Stat(payloadJsonPath + ".jwt"); os.IsNotExist(err) {
		tokenStatus = "created"
	} else {
		tokenStatus = "updated"
	}

	err = ioutil.WriteFile(payloadJsonPath+".jwt", []byte(tokenString), 0640)
	if err != nil {
		return "", err
	}

	log.Infof("Token %s %s.", payloadJsonPath+".jwt", tokenStatus)

	return tokenString, nil
}

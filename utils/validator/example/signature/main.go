package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

type Parameter struct {
	payload    string
	signingKey string
	signature  string
}

func main() {
	payload := `{"transaction_id":"trx001","issuer_product_id":"TSEL10K","customer_number":"081234567890","issuer_code":"dummy"}`
	secret := `123456`

	param := new(Parameter)
	param.payload = payload
	param.signingKey = secret

	if err := bindSignature(param); err != nil {
		fmt.Println("ERROR: ", err.Error())
	} else {
		fmt.Println("Signature:")
		fmt.Println(param.signature)
	}
}

func bindSignature(param *Parameter) (err error) {
	reqTime := time.Now().Unix()
	payloadMinified, err := JSONstringify([]byte(param.payload))
	if err != nil {
		return
	}

	reqTimeString := strconv.FormatInt(reqTime, 10)
	tokenString := reqTimeString + ":" + string(payloadMinified)
	h := hmac.New(sha256.New, []byte(param.signingKey))
	h.Write([]byte(tokenString))
	tokenHashed := hex.EncodeToString(h.Sum(nil))

	param.signature = base64.URLEncoding.EncodeToString([]byte(reqTimeString + ":" + tokenHashed))

	return
}

func JSONstringify(data []byte) ([]byte, error) {
	buff := new(bytes.Buffer)
	errCompact := json.Compact(buff, data)
	if errCompact != nil {
		newErr := fmt.Errorf("failure encountered compacting json := %v", errCompact)
		return nil, newErr
	}

	b, err := ioutil.ReadAll(buff)
	if err != nil {
		readErr := fmt.Errorf("read buffer error encountered := %v", err)
		return nil, readErr
	}

	return b, nil
}

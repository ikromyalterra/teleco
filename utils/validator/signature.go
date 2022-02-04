package validator

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
	"time"
)

type (
	Validator struct {
		ReqTime    string
		ReqPayload string
		Err        error
		IsValid    bool
		Checked    bool
		signature  Signature
	}

	Signature struct {
		Token     string
		Secret    string
		Payload   []byte
		TimeLimit int
	}
)

var (
	ErrToParseKey       = "unable to parse key"
	ErrKeyFormat        = "invalid key format"
	ErrToParseTimestamp = "unable to parse timestamp"
)

func NewSignatureValidator(signature Signature) Validator {
	return Validator{
		signature: signature,
	}
}

func (v Validator) Bind(f func(val Validator) Validator) Validator {
	var checkedButInvalid = v.Checked && !v.IsValid
	if v.Err != nil || checkedButInvalid {
		return v
	}
	return f(v)
}

func (v Validator) Verify() (bool, error) {
	if !v.Checked {
		v.IsValid = false
	}
	return v.IsValid, v.Err
}

func Parse(v Validator) Validator {
	decodedToken, err := base64.URLEncoding.DecodeString(v.signature.Token)
	if err != nil {
		v.Err = errors.New(ErrToParseKey)
		return v
	}

	strDecodedToken := string(decodedToken)
	datas := strings.Split(strDecodedToken, ":")
	if len(datas) != 2 {
		v.Err = errors.New(ErrKeyFormat)
		return v
	}

	v.ReqTime = datas[0]
	v.ReqPayload = datas[1]

	return v
}

func VerifyReqTime(v Validator) Validator {
	reqTimestamp, err := strconv.ParseInt(v.ReqTime, 10, 64)
	if err != nil || reqTimestamp == 0 {
		v.Err = errors.New(ErrToParseTimestamp)
		return v
	}

	now := time.Now().Unix()
	diff := now - reqTimestamp
	intDiff := int(diff)

	v.IsValid = intDiff > 0 && v.signature.TimeLimit >= intDiff
	v.Checked = true

	return v
}

func VerifyReqPayload(v Validator) Validator {
	digest := hmac.New(sha256.New, []byte(v.signature.Secret))
	payload := v.ReqTime + ":" + string(v.signature.Payload)
	digest.Write([]byte(payload))
	expectedHashedPayload := hex.EncodeToString(digest.Sum(nil))

	v.IsValid = hmac.Equal([]byte(v.ReqPayload), []byte(expectedHashedPayload))
	v.Checked = true

	return v
}

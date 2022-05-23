package webapi

import (
	"encoding/base64"
	"errors"
	"strings"
)

// Currently not used
func jwtVerify(jwt []byte) (bool, error) {
	var decoded []byte
	_, err := base64.RawURLEncoding.Decode(jwt, decoded)
	if err != nil {
		return false, err
	}
	j3t := strings.Split(string(decoded), ".")
	if len(j3t) != 3 {
		return false, errors.New("验证字段长度有误")
	}

	return true, nil
}

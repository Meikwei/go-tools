package encoding

import (
	"encoding/base64"

	"github.com/Meikwei/go-tools/errs"
)


func Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}


func Base64Decode(data string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errs.WrapMsg(err, "DecodeString failed", "data", data)
	}
	return string(decodedBytes), nil
}

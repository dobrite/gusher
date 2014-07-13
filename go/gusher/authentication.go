package gusher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

var key = "278d425bdf160c739803"
var secret = "7ad3773142a6692b25b8"

func auth(data ...string) string {
	joined := strings.Join(data, ":")
	return key + ":" + sign(joined)
}

func sign(data string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

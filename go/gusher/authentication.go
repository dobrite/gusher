package gusher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"strings"
)

var keyFlag = flag.String("key", "", "auth key")
var secretFlag = flag.String("secret", "", "auth secret")

func auth(data ...string) string {
	joined := strings.Join(data, ":")

	// TODO move to config object
	// error out if key ends up empty
	var key string

	key = os.Getenv("GUSHER_KEY")

	if key == "" {
		key = *keyFlag
	}

	return key + ":" + sign(joined)
}

func sign(data string) string {
	// TODO move to config object
	// error out if secret ends up empty
	var secret string

	secret = os.Getenv("GUSHER_SECRET")

	if secret == "" {
		secret = *secretFlag
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

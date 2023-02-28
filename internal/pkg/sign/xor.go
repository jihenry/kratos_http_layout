package sign

import (
	"encoding/base64"
)

type rsaConfig struct {
	PrivateKey []byte
	PublicKey  []byte
	XOR        string
}

var (
	blockMap = map[string]*rsaConfig{
		"pocket": {
			XOR: "dnJAPa3ndA1JLsuY",
		},
	}
)

func XOREncode(msg, key string) string {
	ml := len(msg)
	XOR := blockMap[key].XOR
	kl := len(blockMap[key].XOR)
	pwd := ""
	for i := 0; i < ml; i++ {
		pwd += string((XOR[i%kl]) ^ (msg[i]))
	}
	return base64.StdEncoding.EncodeToString([]byte(pwd))
}

func XORDecode(msg, key string) string {
	str, _ := base64.StdEncoding.DecodeString(msg)
	msg = string(str)
	XOR := blockMap[key].XOR
	ml := len(msg)
	kl := len(XOR)
	pwd := ""
	for i := 0; i < ml; i++ {
		pwd += string((msg[i]) ^ XOR[i%kl])
	}
	return pwd
}

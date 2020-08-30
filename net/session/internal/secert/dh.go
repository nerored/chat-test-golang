/*
	密钥交换
*/
package secert

import (
	"time"

	"github.com/nerored/chat-test-golang/net/socket/utils"
	"golang.org/x/crypto/curve25519"
)

const (
	ECDH_KEY_LEN = 32
)

func ECDHKeyNew() (publicKey, privateKey []byte, err error) {
	privateKey = utils.NewMT19937(uint64(time.Now().UnixNano())).Rand32BytesKey()
	publicKey, err = curve25519.X25519(privateKey, curve25519.Basepoint)
	return
}

func ECDHKeyGen(publicKey, privateKey []byte) (key []byte, err error) {
	return curve25519.X25519(privateKey, publicKey)
}

package secert

import (
	"bytes"
	"testing"
)

func TestCryptAES256(t *testing.T) {
	publicKey1, privateKey1, _ := ECDHKeyNew()
	publicKey2, privateKey2, _ := ECDHKeyNew()

	secret1, err := ECDHKeyGen(publicKey1, privateKey2)

	if err != nil {
		t.Fatalf("gen secret1 failed %v", err)
	}

	secret2, err := ECDHKeyGen(publicKey2, privateKey1)

	if err != nil {
		t.Fatalf("gen secret2 failed %v", err)
	}

	if !bytes.Equal(secret1[:], secret2[:]) {
		t.Fatalf("key switch failed\n1\t%v\n2\t%v", secret1, secret2)
	}

	aesBlock, err := NewAESBlockCrypt(secret1[:])

	if aesBlock == nil || err != nil {
		t.Fatalf("create block failed err %v", err)
	}

	sourcedata := []byte("hello ase 256")

	encodeData := make([]byte, len(sourcedata))
	decodeData := make([]byte, len(sourcedata))

	aesBlock.Encrypt(encodeData, sourcedata)
	aesBlock.Decrypt(decodeData, encodeData)

	if !bytes.Equal(sourcedata, decodeData) {
		t.Fatalf("decode failed")
	}

	t.Logf("src %+v\nencode %+v\ndecode %+v", sourcedata, encodeData, decodeData)
}

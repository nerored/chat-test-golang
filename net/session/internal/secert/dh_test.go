package secert

import (
	"bytes"
	"testing"
)

func TestKeySwitch(t *testing.T) {
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

	t.Logf("publicKey1 %v privateKey1 %v secret1 %v", publicKey1, privateKey1, secret1)
	t.Logf("publicKey2 %v privateKey2 %v secret2 %v", publicKey2, privateKey2, secret2)
}

package encrypt

import (
	"testing"
)

func TestEncoding(t *testing.T) {
	secret := "1234567890123456"
	encrypted, ee := AesEncrypt("i love u", secret)
	if ee != nil {
		t.Error(ee)
	}

	decrypted, de := AesDecrypt(encrypted, secret)
	if de != nil {
		t.Error(de)
	}

	if decrypted != "i love u" {
		t.Error("decrypted message not match")
	}

	t.Log("encrypted message:", encrypted, "decrypted message:", decrypted)
}

func TestPasswd(t *testing.T) {
	encoded, ee := BcryptEncode("i love u")
	if ee != nil {
		t.Error(ee)
	}

	if !BcryptCheck("i love u", encoded) {
		t.Error("bcrypt check failed")
	}

	t.Log("encoded message:", encoded)
}

type entry struct {
	Name string
	Age  int
}

func TestHashMD5(t *testing.T) {
	t.Run("HashMD5", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Log("hash md5:", HashMD5("i love u"))
		}
	})

	t.Run("HashEntryMD5", func(t *testing.T) {
		entry := &entry{
			Name: "alice",
			Age:  114,
		}

		for i := 0; i < 10; i++ {
			t.Log("hash md5:", HashEntryMD5(entry))
		}
	})
}

func TestRsa(t *testing.T) {
	t.Run("RsaEncrypt:Success", func(t *testing.T) {
		pri, pub, err := RsaKeyGenerate(2048)
		if err != nil {
			t.Error(err)
		}

		t.Log("private key:", pri)
		t.Log("public key:", pub)

		encrypted, ee := RsaEncrypt(pub, "i love u")
		if ee != nil {
			t.Error(ee)
		}

		decrypted, de := RsaDecrypt(pri, encrypted)
		if de != nil {
			t.Error(de)
		}

		if decrypted != "i love u" {
			t.Error("decrypted message not match")
		}
	})

	t.Run("RsaEncrypt:NoKey", func(t *testing.T) {
		pri, pub, err := RsaKeyGenerate(2048)
		if err != nil {
			t.Error(err)
		}

		t.Log("private key:", pri)
		t.Log("public key:", pub)

		encrypted, ee := RsaEncrypt("", "i love u")
		if ee == nil {
			t.Error("want error, but got nil")
		}

		_, de := RsaDecrypt("", encrypted)
		if de == nil {
			t.Error("want error, but got nil")
		}
	})
}

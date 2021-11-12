package crypto_test

import (
	"testing"

	"github.com/sisukasco/commons/crypto"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestSignatureVarification(t *testing.T) {
	key := crypto.GeneratePKCKeys()

	for i := 0; i < 10; i++ {
		data := struct {
			Value1 string
			Value2 string
			Value3 int64
		}{
			Value1: faker.RandomString(22),
			Value2: faker.RandomString(22),
			Value3: faker.RandomInt64(1000, 1000000),
		}

		sign, err := crypto.GenerateSignature(&data, key.PrivateKey)
		assert.Nil(t, err)

		err = crypto.VerifySignature(&data, key.PublicKey, sign)
		assert.Nil(t, err)
	}

}

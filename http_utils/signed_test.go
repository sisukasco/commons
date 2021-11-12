package http_utils_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/sisukas/commons/crypto"
	"github.com/sisukas/commons/http_utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

type SimpleStruct struct {
	Value1 string
	Value2 string
}

func TestPostingSignedRequest(t *testing.T) {
	key := crypto.GeneratePKCKeys()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		signx, err := http_utils.ExtractSignature(r)
		assert.Nil(t, err)
		t.Logf("signature %s", signx)
		ss := &SimpleStruct{}

		jsonDecoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		err = jsonDecoder.Decode(ss)
		assert.Nil(t, err)

		err = crypto.VerifySignature(ss, key.PublicKey, signx)
		assert.Nil(t, err)

		t.Logf("Verified")
	}))
	defer svr.Close()

	ctx := context.Background()

	body := &SimpleStruct{
		Value1: faker.RandomString(12),
		Value2: faker.RandomString(22),
	}
	_, err := http_utils.PostSignedRequest(ctx, svr.URL, body, key.PrivateKey)
	assert.Nil(t, err)

}

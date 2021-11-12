package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
)

//ED25519

func GenerateSignature(msg interface{}, privateKeyB64 string) (string, error) {
	privKey, err := base64.StdEncoding.DecodeString(privateKeyB64)
	if err != nil {
		return "", err
	}
	bMsg, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	sign := ed25519.Sign(privKey, bMsg)

	return base64.URLEncoding.EncodeToString(sign), nil
}

func VerifySignature(msg interface{}, publicKeyB64 string, strSign string) error {
	pubKey, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return err
	}
	sign, err := base64.URLEncoding.DecodeString(strSign)
	if err != nil {
		return err
	}
	bMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	bVerified := ed25519.Verify(pubKey, bMsg, sign)
	if !bVerified {
		return errors.New("Signature does not match")
	}
	return nil
}

type PKCKeyCombo struct {
	PublicKey  string
	PrivateKey string
}

func GeneratePKCKeys() *PKCKeyCombo {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pkkey := &PKCKeyCombo{
		PublicKey:  base64.StdEncoding.EncodeToString(publicKey),
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
	}

	/*fmt.Printf("Public Key:\n %s \nPrivateKey:\n%s\n",
	base64.StdEncoding.EncodeToString(publicKey),
	base64.StdEncoding.EncodeToString(privateKey))
	*/
	return pkkey
}

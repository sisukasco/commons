package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io"
)

func pkcs7strip(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("pkcs7: Data is empty")
	}
	if length%blockSize != 0 {
		return nil, errors.New("pkcs7: Data is not block-aligned")
	}
	padLen := int(data[length-1])
	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(data, ref) {
		return nil, errors.New("pkcs7: Invalid padding")
	}
	return data[:length-padLen], nil
}

// pkcs7pad add pkcs7 padding
func pkcs7pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize < 0 || blockSize > 256 {
		return nil, errors.New(fmt.Sprintf("pkcs7: Invalid block size %d", blockSize))
	} else {
		padLen := blockSize - len(data)%blockSize
		padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
		return append(data, padding...), nil
	}
}

func DecryptAES(c string, k string) (string, error) {

	bRawData, err := base64.RawURLEncoding.DecodeString(c)
	if err != nil {
		return "", errors.Wrap(err, "Error decrypting")
	}
	iv := bRawData[:16]
	cs := bRawData[16:]

	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		return "", errors.Wrap(err, "error creating cipher block")
	}
	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(cs, cs)

	strp, err := pkcs7strip(cs, 32)
	if err != nil {
		return "", errors.Wrap(err, "Error stripping padding")
	}

	return string(strp), nil
}

func URLEncodedBase64ToHex(c string) (string, error) {

	bRawData, err := base64.RawURLEncoding.DecodeString(c)
	if err != nil {
		return "", errors.Wrap(err, "Error decoding")
	}
	return hex.EncodeToString(bRawData), nil
}

func EncryptAES(plain string, key string) (string, error) {

	plaintext, err := pkcs7pad([]byte(plain), 32)
	if err != nil {
		return "", errors.Wrapf(err, "Can't padd the plain text ")
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.Wrapf(err, "Can't make iv")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", errors.Wrapf(err, "Can't make NewCipher")
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	encr := base64.RawURLEncoding.EncodeToString(ciphertext)

	return encr, nil

}

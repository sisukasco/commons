package http_utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/sisukasco/commons/crypto"
)

type Response struct {
	Body string
}

func PostSignedRequest(ctx context.Context, url string, body interface{}, privateKeyB64 string) (*Response, error) {
	bBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	sign, err := crypto.GenerateSignature(body, privateKeyB64)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost,
		url, bytes.NewReader(bBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	request.Header.Set("Authorization", fmt.Sprintf("Signature %s", sign))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Received response was not 200")
	}
	bytesBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &Response{Body: string(bytesBody)}, nil
}

var signatureRegexp = regexp.MustCompile(`^(?:S|s)ignature (\S+$)`)

func ExtractSignature(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	return parseSignature(authHeader)
}

func parseSignature(authHeader string) (string, error) {
	if len(authHeader) <= 0 {
		return "", errors.New("Signature is empty")
	}
	matches := signatureRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("Couldn't extract signature from header")
	}
	if len(matches[1]) <= 5 {
		return "", errors.New("Couldn't extract signature from header")
	}
	return matches[1], nil
}

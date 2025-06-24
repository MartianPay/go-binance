package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Signer struct {
	APIKey    string
	SecretKey string
}

func NewSigner(apiKey, secretKey string) *Signer {
	return &Signer{
		APIKey:    apiKey,
		SecretKey: secretKey,
	}
}

func (s *Signer) Sign(params map[string]string) string {
	if params == nil {
		params = make(map[string]string)
	}

	params["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)

	queryString := s.BuildQueryString(params)
	signature := s.GenerateSignature(queryString)

	return queryString + "&signature=" + signature
}

func (s *Signer) BuildQueryString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var pairs []string
	for _, k := range keys {
		if v := params[k]; v != "" {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		}
	}

	return strings.Join(pairs, "&")
}

func (s *Signer) GenerateSignature(queryString string) string {
	h := hmac.New(sha256.New, []byte(s.SecretKey))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Signer) GetHeaders() map[string]string {
	return map[string]string{
		"X-MBX-APIKEY": s.APIKey,
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	LanguageJa string = "ja"
	LanguageEn string = "en"
)

var (
	langMap = map[string]bool{
		LanguageJa: true,
		LanguageEn: true,
	}

	translateURL = mustParseURL()
)

type (
	APIParams struct {
		Text   string
		Source string
		Target string
	}

	APICaller struct {
		client *http.Client
		params *APIParams
	}

	APIResponce struct {
		Code int    `json:"code"`
		Text string `json:"text"`
	}
)

func mustParseURL() *url.URL {
	b, err := ioutil.ReadFile("url.txt")
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(string(b))
	if err != nil {
		panic(err)
	}

	return u
}

func NewAPICaller(p *APIParams) *APICaller {
	return &APICaller{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		params: p,
	}
}

func (caller APICaller) Call() (string, error) {
	var (
		p = caller.params
		q = translateURL.Query()
	)

	q.Add("text", p.Text)
	q.Add("source", p.Source)
	q.Add("target", p.Target)
	translateURL.RawQuery = q.Encode()

	resp, err := caller.client.Get(translateURL.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	json := &APIResponce{}
	err = decoder.Decode(json)
	if err != nil {
		return "", err
	}

	if json.Code == http.StatusOK {
		return json.Text, nil
	}

	return "", fmt.Errorf("Error Code: %d", http.StatusBadRequest)
}

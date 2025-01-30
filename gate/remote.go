// Here we connect to a remote microservice

package main

import (
	"encoding/json"
	"fmt"
	"gate/api"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type NipData struct {
	ShortName string `json:"shortName"`
	LongName  string `json:"longName"`
	TaxId     string `json:"taxId"`
	Apartment string `json:"apartment"`
	Building  string `json:"building"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Zip       string `json:"zip"`
}

func getNipData(nip string) (*api.NipResponse, error) {

	remoteUrl := "http://185.2.114.92:20011/getcontractor"
	remoteAuthToken := "michalspitoken"

	reqBody := url.Values{}
	reqBody.Set("nip", nip)

	req, err := http.NewRequest(http.MethodPost, remoteUrl, strings.NewReader(reqBody.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("AUTH-TOKEN", remoteAuthToken)

	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var nipResponse api.NipResponse
	if err := json.Unmarshal(body, &nipResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	return &nipResponse, nil
}

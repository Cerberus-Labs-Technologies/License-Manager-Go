package rest

import (
	"CerberusLabsLicenseManagerGo/license"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RestResult struct {
	Data    license.License `json:"data"`
	Success bool            `json:"success"`
}

func MakeHttpGetRequest(url string) (license.License, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return license.License{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return license.License{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 403:
			return license.License{}, fmt.Errorf("The license is inactive!")
		case 401:
			return license.License{}, fmt.Errorf("The license is blocked!")
		case 406:
			return license.License{}, fmt.Errorf("The license is expired!")
		case 404:
			return license.License{}, fmt.Errorf("The license is not found!")
		default:
			return license.License{}, fmt.Errorf("An error occurred while trying to connect to the license server!")
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return license.License{}, err
	}

	restResult := RestResult{}
	err = json.Unmarshal(body, &restResult)
	if err != nil {
		return license.License{}, err
	}

	if !restResult.Success {
		return license.License{}, fmt.Errorf("An error occurred while trying to connect to the license server!")
	}

	return restResult.Data, nil
}

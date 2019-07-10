package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RatesApiService struct {
}

const (
	apiUrl = "https://api.ratesapi.io/api/"
)

func (s *RatesApiService) Get(day string) (Rate, error) {
	httpClient := http.Client{}

	callURL := apiUrl
	if day == "" {
		callURL += "latest"
	} else {
		callURL += day
	}

	request, err := http.NewRequest(http.MethodGet, callURL, nil)

	if err != nil {
		return Rate{}, err
	}

	res, err := httpClient.Do(request)
	if err != nil {
		return Rate{}, err
	}

	if res.StatusCode != http.StatusOK {
		return Rate{}, errors.New("bad status code " + res.Status)
	}

	var rate Rate
	if err := json.NewDecoder(res.Body).Decode(&rate); err != nil {
		return Rate{}, err
	}
	return rate, nil

}

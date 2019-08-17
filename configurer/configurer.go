package configurer

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	httpclient "dat.havinh/week5-demo/http-client"
)

const (
	activatedServiceIDsEndpoint      = "http://my-json-server.typicode.com/havinhdat/restful-db/config/activatedServiceIDs"
	enabledSubTextServiceIDsEndpoint = "http://my-json-server.typicode.com/havinhdat/restful-db/config/allowedSubTextServiceIDs"
)

//go:generate mockery -name=Configurer
type Configurer interface {
	GetActivatedServiceIDs() ([]int64, error)
	GetEnabledSubTextServiceIDs() ([]int64, error)
}

type configurerImpl struct {
	httpClient httpclient.HTTPClient
}

type ListServiceIDsConfig struct {
	ID    string  `json:"id"`
	Value []int64 `json:"value"`
}

func (cfg *configurerImpl) GetActivatedServiceIDs() ([]int64, error) {
	return cfg.getListServiceIDsConfig(activatedServiceIDsEndpoint)
}

func (cfg *configurerImpl) GetEnabledSubTextServiceIDs() ([]int64, error) {
	return cfg.getListServiceIDsConfig(enabledSubTextServiceIDsEndpoint)
}

func (cfg *configurerImpl) getListServiceIDsConfig(configEndpoint string) ([]int64, error) {
	resp, err := cfg.httpClient.Get(configEndpoint)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()
	configVal := &ListServiceIDsConfig{}
	if err = json.Unmarshal(body, configVal); err != nil {
		return nil, err
	}
	return configVal.Value, nil
}

type Option func(configurer *configurerImpl)

func New(options ...Option) (Configurer, error) {
	configurer := &configurerImpl{}
	for _, o := range options {
		o(configurer)
	}
	if configurer.httpClient == nil {
		return nil, errors.New("missing http client")
	}
	return configurer, nil
}

func WithHTTPClient(httpClient httpclient.HTTPClient) Option {
	return func(configurer *configurerImpl) {
		configurer.httpClient = httpClient
	}
}

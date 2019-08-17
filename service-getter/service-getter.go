package service_getter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	httpclient "dat.havinh/week5-demo/http-client"
)

const (
	getServiceEndpoint = "http://my-json-server.typicode.com/havinhdat/restful-db/services/%d"
)

type ServiceGetter interface {
	Get(id int64) (*Service, error)
}

type Service struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	SubText string `json:"subText"`
}

func (s Service) Format(withSubtext bool) string {
	subText := s.SubText
	if withSubtext {
		subText = "<omitted>"
	}
	return fmt.Sprintf("Service: %s\nSubText: %s", s.Name, subText)
}

type serviceGetterImpl struct {
	httpClient httpclient.HTTPClient
}

func (sg *serviceGetterImpl) Get(id int64) (*Service, error) {
	resp, err := sg.httpClient.Get(fmt.Sprintf(getServiceEndpoint, id))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer func() {
		_ = resp.Body.Close()
	}()

	service := &Service{}
	if err := json.Unmarshal(body, service); err != nil {
		return nil, err
	}

	return service, nil
}

func New(httpClient httpclient.HTTPClient) (ServiceGetter, error) {
	return &serviceGetterImpl{
		httpClient: httpClient,
	}, nil
}

package configurer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	httpclient "dat.havinh/week5-demo/http-client"
)

func TestTestGetActivatedServiceIDs_New(t *testing.T) {
	var (
		expectedServiceIDs = []int64{1, 5, 7}
		nilServiceIDs      []int64
	)
	testCases := []struct {
		desc               string
		doMock             func(httpClientMocks *httpclient.MockHTTPClient)
		expectedServiceIDs []int64
		expectedError      error
	}{
		{
			desc:               "success",
			expectedServiceIDs: expectedServiceIDs,
			expectedError:      nil,
			doMock: func(httpClientMocks *httpclient.MockHTTPClient) {
				listMocks := &ListServiceIDsConfig{
					ID:    "anything",
					Value: expectedServiceIDs,
				}
				bytes, _ := json.Marshal(listMocks)
				respMocks := &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(string(bytes))),
				}
				httpClientMocks.On("Get", activatedServiceIDsEndpoint).Return(respMocks, nil)
			},
		},
		{
			desc:               "failed",
			expectedServiceIDs: nilServiceIDs,
			expectedError:      errors.New("anything"),
			doMock: func(httpClientMocks *httpclient.MockHTTPClient) {
				httpClientMocks.On("Get", activatedServiceIDsEndpoint).Return(nil, errors.New("anything"))
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			httpClientMocks := &httpclient.MockHTTPClient{}
			conf, _ := New(WithHTTPClient(httpClientMocks))

			if tc.doMock != nil {
				tc.doMock(httpClientMocks)
			}

			serviceIDs, err := conf.GetActivatedServiceIDs()

			assert.Equal(t, tc.expectedServiceIDs, serviceIDs)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetActivatedServiceIDs(t *testing.T) {
	httpClientMocks := &httpclient.MockHTTPClient{}
	conf, _ := New(WithHTTPClient(httpClientMocks))

	serviceIDsMocks := []int64{1, 5, 7}
	listMocks := &ListServiceIDsConfig{
		ID:    "anything",
		Value: serviceIDsMocks,
	}
	bytes, _ := json.Marshal(listMocks)
	respMocks := &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(string(bytes))),
	}
	httpClientMocks.On("Get", activatedServiceIDsEndpoint).Return(respMocks, nil)

	serviceIDs, err := conf.GetActivatedServiceIDs()
	assert.Equal(t, serviceIDsMocks, serviceIDs)
	assert.Equal(t, nil, err)
}

func TestGetActivatedServiceIDs_Failed(t *testing.T) {
	httpClientMocks := &httpclient.MockHTTPClient{}
	conf, _ := New(WithHTTPClient(httpClientMocks))
	var nilServiceIDs []int64

	httpClientMocks.On("Get", activatedServiceIDsEndpoint).Return(nil, errors.New("failed http"))

	serviceIDs, err := conf.GetActivatedServiceIDs()
	assert.Equal(t, nilServiceIDs, serviceIDs)
	assert.Equal(t, errors.New("failed http"), err)
}

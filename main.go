package main

import (
	"fmt"
	"log"
	"net/http"

	"dat.havinh/week5-demo/configurer"
	service_getter "dat.havinh/week5-demo/service-getter"
)

// Join Facebook Group: http://go.grab.com/vnbootcamp-week5-fb

/*
Requirement:
Show Grab service details
- Get activated service ids : configurer
- Get service details : service_getter
- Check if subtext is required : configurer

GET http://my-json-server.typicode.com/havinhdat/restful-db/services/{id}
{
	"id": $int64,
	"name": $string,
	"subText": $string
}
GET http://my-json-server.typicode.com/havinhdat/restful-db/config/{id}
{
	"id": $string,
	"value": $interface
}
activatedServiceIDs: value is slice of service ids
allowedSubTextServiceIDs: value is slice of service ids
*/

func main() {
	httpClient := http.DefaultClient
	conf, configurerErr := configurer.New(configurer.WithHTTPClient(httpClient))
	serviceGetter, serviceGetterErr := service_getter.New(httpClient)
	if configurerErr != nil || serviceGetterErr != nil {
		log.Println("ERROR: failed to init configurer or service getter")
		return
	}

	resp, err := Handle(conf, serviceGetter)
	if err != nil {
		log.Println("ERROR: failed to handle with error: ", err)
	}
	for _, service := range resp.Services {
		fmt.Printf(service.Format(!resp.EnabledSubText[service.ID]))
		fmt.Printf("\n\n")
	}
}

type Response struct {
	Services       []*service_getter.Service
	EnabledSubText map[int64]bool
}

func Handle(conf configurer.Configurer, serviceGetter service_getter.ServiceGetter) (*Response, error) {
	activatedServiceIDs, err := conf.GetActivatedServiceIDs()
	if err != nil {
		return nil, err
	}
	services := make([]*service_getter.Service, 0, len(activatedServiceIDs))
	for _, serviceID := range activatedServiceIDs {
		s, err := serviceGetter.Get(serviceID)
		if err != nil {
			log.Println("WARN: failed to get service by id with error: ", err)
			continue
		}
		services = append(services, s)
	}
	enabledSubTextServiceIDs, err := conf.GetEnabledSubTextServiceIDs()
	if err != nil {
		return nil, err
	}
	enabledSubTextByServiceID := map[int64]bool{}
	for _, serviceID := range enabledSubTextServiceIDs {
		enabledSubTextByServiceID[serviceID] = true
	}

	return &Response{
		Services:       services,
		EnabledSubText: enabledSubTextByServiceID,
	}, nil
}

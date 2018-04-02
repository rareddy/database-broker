package broker

import (
	"encoding/json"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"io/ioutil"
)

type datasource struct {
	Name        string
	Parameters  map[string]interface{}
}

type externaldatasource interface {
	defaultPort() int
	name() string
}

type configurable interface {
	springboot(bindAlias string, multiSource bool) map[string]interface{}
	wildflyswarm(bindAlias string, multiSource bool) map[string]interface{}
	nodejs(bindAlias string, multiSource bool) map[string]interface{}
	other(bindAlias string, multiSource bool) map[string]interface{}
}

func catalog() ([]osb.Service, error) {
	response := &osb.CatalogResponse{}

	data, err := ioutil.ReadFile("/opt/servicebroker/service.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		return nil, err
	}
	return response.Services, nil
}

func newDataSource(serviceInstance DataSourceInstance, bindingParameters map[string]interface{}) configurable {
	services, err := catalog()
	if err != nil {
		// we should never get here.
		panic(err)
	}

	id := serviceInstance.PlanID
	for i := range services {
		for j := range services[i].Plans {
			if services[i].Plans[j].ID == id {
				ds := datasource{}
				ds.Name = services[i].Name
				ds.Parameters = merge(serviceInstance.Parameters, bindingParameters)
				// TODO: need to figure out if there is a way to reflectively do this?
				if services[i].Plans[j].Name == "mysql" {
					return MySql{ds}
				} else if services[i].Plans[j].Name == "sqlserver" {
					return SqlServer{ds}
				} else if services[i].Plans[j].Name == "oracle" {
					return Oracle{ds}
				} else if services[i].Plans[j].Name == "postgresql" {
					return PostgreSQL{ds}
				}
			}
		}
	}
	panic("did not find the datasource with given id " + id)
}

func merge(properties ...map[string]interface{}) map[string]interface{} {
	creds := make(map[string]interface{})

	for i := range properties {
		for k, v := range properties[i] {
			creds[k] = v
		}
	}
	return creds
}

func wfs(bindAlias string, key string) string {
	return "swarm.datasources.data-sources."+ bindAlias +"."+key
}

func sb(bindAlias string, key string, multiSource bool) string {
	if multiSource {
		return "spring.datasource." + bindAlias + "." + key
	}
	return "spring.datasource." + key
}

func (ds datasource) springboot() map[string]interface{} {
	panic("should not have reached")
}

func (ds datasource) wildflyswarm() map[string]interface{} {
	panic("should not have reached")}

func (ds datasource) nodejs() map[string]interface{} {
	panic("should not have reached")}

func (ds datasource) other() map[string]interface{} {
	panic("should not have reached")
}
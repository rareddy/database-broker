package broker

import (
	"encoding/json"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"io/ioutil"
)

type datasource struct {
	Name              string
	SourceName        string
	Parameters map[string]interface{}
}

type externaldatasource interface {
	defaultPort() int
}

type configurable interface {
	springboot() map[string]interface{}
	wildflyswarm() map[string]interface{}
	nodejs() map[string]interface{}
	other() map[string]interface{}
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

func newDataSource(serviceInstance dbServiceInstance, bindingParameters map[string]interface{}) configurable {
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
				ds.SourceName = i2s(serviceInstance.Parameters["source-name"])
				ds.Parameters = merge(serviceInstance.Parameters, bindingParameters)
				// TODO: need to figure out if there is a way to reflectively do this?
				if services[i].Plans[j].Name == "mysql" {
					return MySql{ds}
				} else if services[i].Plans[j].Name == "sqlserver" {
					return SqlServer{ds}
				} else if services[i].Plans[j].Name == "oracle" {
					return Oracle{ds}
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

func wfs(sourceName string, key string) string {
	return "swarm.datasources.data-sources."+sourceName+"."+key
}

func sb(sourceName string, key string) string {
	return "spring.datasource."+sourceName+"."+key
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
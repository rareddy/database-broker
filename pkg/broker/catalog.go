package broker

import (
	"encoding/json"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"io/ioutil"
	"path/filepath"
)
/*
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
*/
func catalog() ([]osb.Service, error) {

	databaseBroker := osb.Service{}

	databaseBroker.Name = "database-broker"
	databaseBroker.ID = "5b333b4a-e37c-4391-9090-af59dfcaa1cc"
	databaseBroker.Description = "Provision a Connection to database like Oracle, MS-SQLServer, MySQL on your network."
	databaseBroker.Tags = []string{"no-sql", "database"}
	databaseBroker.Bindable = true

	metadata := map[string]interface{}{}
	metadata["displayName"] = "Database Broker"
	metadata["imageUrl"] = "https://avatars2.githubusercontent.com/u/19862012?s=200&v=4"
	provider := map[string]interface{}{}
	provider["name"] = "Database Broker"
	metadata["provider"] = provider

	databaseBroker.Metadata = metadata
	planUpdatable := false
	databaseBroker.PlanUpdatable = &planUpdatable

	files, err := filepath.Glob("/opt/servicebroker/plans/*.json")
	if err == nil {
		for _, f := range files {
			data, err := ioutil.ReadFile(f)
			if err == nil {
				plan := osb.Plan{}
				err = json.Unmarshal([]byte(data), &plan)
				if err == nil {
					databaseBroker.Plans = append(databaseBroker.Plans, plan)
				}
			}
		}
	}
	return []osb.Service{databaseBroker}, nil
}
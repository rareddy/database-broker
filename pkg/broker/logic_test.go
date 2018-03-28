package broker

import (
	"fmt"
	"github.com/pmorie/osb-broker-lib/pkg/broker"
	"testing"
	"encoding/json"
)

func TestJSON(t *testing.T) {
	var request broker.RequestContext
	var logic BusinessLogic
	response, err := logic.GetCatalog(&request)
	if err != nil {
		fmt.Println(err)
	}
	if response.Services[0].ID != "5b333b4a-e37c-4391-9090-af59dfcaa1cc" {
		t.Error("wrong id received")	
	}
	
	fmt.Println("json ")
	b, err := json.MarshalIndent(response, "", "  ");
	fmt.Println(string(b))
}

package broker

import (
	"net/http"
	"sync"

	"github.com/pmorie/osb-broker-lib/pkg/broker"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"github.com/golang/glog"
	"errors"
)

// NewBusinessLogic is a hook that is called with the Options the program is run
// with. NewBusinessLogic is the place where you will initialize your
// BusinessLogic the parameters passed in.
func NewBusinessLogic(o Options) (*BusinessLogic, error) {
	// For example, if your BusinessLogic requires a parameter from the command
	// line, you would unpack it from the Options and set it on the
	// BusinessLogic here.
	return &BusinessLogic{
		async:     o.Async,
		instances: make(map[string]*dbServiceInstance, 10),
	}, nil
}

// BusinessLogic provides an implementation of the broker.BusinessLogic
// interface.
type BusinessLogic struct {
	// Indicates if the broker should handle the requests asynchronously.
	async bool
	// Synchronize go routines.
	sync.RWMutex
	// Add fields here! These fields are provided purely as an example
	instances map[string]*dbServiceInstance
}

var _ broker.Interface = &BusinessLogic{}

func (b *BusinessLogic) GetCatalog(c *broker.RequestContext) (*osb.CatalogResponse, error) {
	response := &osb.CatalogResponse{}

	services, err := catalog()
	if err != nil {
		return nil, err
	}

	response.Services = services

	return response, nil
}

func (b *BusinessLogic) Provision(request *osb.ProvisionRequest, c *broker.RequestContext) (*osb.ProvisionResponse, error) {

	b.Lock()
	defer b.Unlock()

	// only accept async
	if !request.AcceptsIncomplete {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusUnprocessableEntity,
		}
	}

	// if provision request on existing instance comes in then return 202
	if b.instances[request.InstanceID] != nil {
		glog.Infof("already provisioned ", request.InstanceID)
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusAccepted,
		}
	}

	glog.Infof("received provision request for ", request.InstanceID)

	response := osb.ProvisionResponse{}
	operation := osb.OperationKey("provision")
	serviceInstance := dbServiceInstance{ID: request.InstanceID, Parameters: request.Parameters, PlanID: request.PlanID,
		State: osb.StateInProgress, OperationKey:operation}
	b.instances[request.InstanceID] = &serviceInstance

	response.Async = true

	// create external service in async mode
	go createService(&serviceInstance)

	url := "http://goolge.com"
	response.DashboardURL = &url
	// bug in OpenShift it is not returning this key back with lastoperation.
	response.OperationKey = &operation
	return &response, nil
}

func (b *BusinessLogic) Deprovision(request *osb.DeprovisionRequest, c *broker.RequestContext) (*osb.DeprovisionResponse, error) {
	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, errors.New("Service not found to deprovision"+request.InstanceID)
	}

	// removing the service
	removeExternalService(*instance)

	b.Lock()
	defer b.Unlock()

	response := osb.DeprovisionResponse{}

	delete(b.instances, request.InstanceID)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) LastOperation(request *osb.LastOperationRequest, c *broker.RequestContext) (*osb.LastOperationResponse, error) {
	serviceInstance := b.instances[request.InstanceID]
	glog.Infof("last operation request on instance", request)

	// TODO:Bug in OpenShift OperationKey is always passed as nil.
	if request.OperationKey == nil {
		request.OperationKey = &serviceInstance.OperationKey
	}

	if *request.OperationKey == osb.OperationKey("provision") {
		glog.Infof("Service has been requested to provision", request.InstanceID)
		response := osb.LastOperationResponse{}
		response.State = serviceInstance.State
		glog.Infof("LastOperation response", response)
		return &response, nil
	}
	return nil, nil
}

func createService(serviceInstance *dbServiceInstance) {
	glog.Infof("starting to create a service for instance", serviceInstance.ID)
	serviceInstance.Lock()
	defer serviceInstance.Unlock()

	serviceInstance.State = osb.StateInProgress
	ok, err := createExternalService(*serviceInstance)
	if !ok {
		serviceInstance.State = osb.StateFailed
		glog.Infof("failed to create external service", err)
	} else {
		// modify the status of the instance
		serviceInstance.State = osb.StateSucceeded
		glog.Infof("done creating a service for instance", serviceInstance.ID)
	}
}

func (b *BusinessLogic) Bind(request *osb.BindRequest, c *broker.RequestContext) (*osb.BindResponse, error) {
	// Your bind business logic goes here

	// example implementation:
	b.Lock()
	defer b.Unlock()

	instance, ok := b.instances[request.InstanceID]
	if !ok {
		return nil, osb.HTTPStatusCodeError{
			StatusCode: http.StatusNotFound,
		}
	}

	creds := createBindingParameters(*instance, request.Parameters)
	response := osb.BindResponse{
		Credentials: creds,
	}

	// asynchronously create PodPreset
	go createPodPreset(*instance, creds)

	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) Unbind(request *osb.UnbindRequest, c *broker.RequestContext) (*osb.UnbindResponse, error) {
	// Your unbind business logic goes here
	return &osb.UnbindResponse{}, nil
}

func (b *BusinessLogic) Update(request *osb.UpdateInstanceRequest, c *broker.RequestContext) (*osb.UpdateInstanceResponse, error) {
	// Your logic for updating a service goes here.
	response := osb.UpdateInstanceResponse{}
	if request.AcceptsIncomplete {
		response.Async = b.async
	}

	return &response, nil
}

func (b *BusinessLogic) ValidateBrokerAPIVersion(version string) error {
	return nil
}

// example types

// exampleInstance is intended as an example of a type that holds information about a service instance
type dbServiceInstance struct {
	ID           string
	Parameters   map[string]interface{}
	State        osb.LastOperationState
	sync.RWMutex
	PlanID       string
	OperationKey osb.OperationKey
}

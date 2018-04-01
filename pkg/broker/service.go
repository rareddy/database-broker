package broker

import (
	"errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"regexp"
	"strconv"
	"github.com/golang/glog"
	"fmt"
)

func serviceAction(serviceInstance dbServiceInstance, opKey osb.OperationKey,
	task func (serviceInstance dbServiceInstance) error,
	after func(osb.OperationKey, error)) {

	glog.Infof("starting to create a service for instance", serviceInstance.ID)
	serviceInstance.Lock()
	defer serviceInstance.Unlock()

	err := task(serviceInstance)
	after(opKey, err)
}

func createExternalService(serviceInstance dbServiceInstance) error {

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	db := newDataSource(serviceInstance, nil)
	serviceName := i2s(serviceInstance.Parameters["source-name"])
	namespace := i2s(serviceInstance.Parameters["namespace"])
	host := i2s(serviceInstance.Parameters["host"])
	portInt64, err := strconv.ParseInt(i2s(serviceInstance.Parameters["port"]), 10, 32)
	if err != nil {
		return errors.New("invalid Port provided")
	}

	port := int32(portInt64)

	service := v1.Service{}
	service.ObjectMeta.Name = serviceName
	service.Spec.Selector = map[string]string{}

	if isIP(host) {
		// if the host matches to IPV4 address
		servicePort := v1.ServicePort{Name: db.(externaldatasource).name(), Port: int32(db.(externaldatasource).defaultPort()),
			Protocol: "TCP", TargetPort: intstr.FromInt(db.(externaldatasource).defaultPort())}
		service.Spec.Ports = []v1.ServicePort{servicePort}
	} else if isHostName(host) {
		// then it assumed this is DNS name used as external name
		service.Spec.Type = "ExternalName"
		service.Spec.ExternalName = host
	} else {
		return errors.New("invalid Hostname provided")
	}

	typeMetadata := metav1.TypeMeta{"Service", "v1"}
	existing, err := clientset.CoreV1().Services(namespace).Get(serviceName, metav1.GetOptions{TypeMeta: typeMetadata})
	if err == nil{
		return errors.New(serviceName + "service already exists, can not provision "+ i2s(existing))
	}

	fmt.Println("no previous service found, creating new...")
	created, err := clientset.CoreV1().Services(namespace).Create(&service)
	if err != nil {
		return err
	}

	glog.Infof("created external service:", created)

	if isIP(host) {
		endpoint := buildEndpoint(serviceName, db.(externaldatasource).name(), host, port)
		createdEndpoint, err := clientset.CoreV1().Endpoints(namespace).Create(&endpoint)
		if err != nil {
			return err
		}
		glog.Infof("created endpoint :", createdEndpoint)
	}
	return nil
}

func buildEndpoint(serviceName string, dbType string, dbHost string, dbPort int32) v1.Endpoints {
	endpoint := v1.Endpoints{}
	endpoint.ObjectMeta.Name = serviceName

	subset := v1.EndpointSubset{}
	address := v1.EndpointAddress{}
	address.IP = dbHost
	subset.Addresses = []v1.EndpointAddress{address}

	endpointPort := v1.EndpointPort{Name: dbType, Port: dbPort}
	subset.Ports = []v1.EndpointPort{endpointPort}

	endpoint.Subsets = []v1.EndpointSubset{subset}
	return endpoint
}

func removeExternalService(serviceInstance dbServiceInstance) (error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	serviceName := i2s(serviceInstance.Parameters["source-name"])
	namespace := i2s(serviceInstance.Parameters["namespace"])
	host := i2s(serviceInstance.Parameters["host"])

	if isIP(host){
		glog.Infof("removing service :", serviceName)
		err = clientset.CoreV1().Services(namespace).Delete(serviceName, nil)
		if err != nil {
			glog.Infof("removing endpoint :", serviceName)
			err = clientset.CoreV1().Endpoints(namespace).Delete(serviceName, nil)
		}
	} else if isHostName(host) {
		glog.Infof("removing service :", serviceName)
		err = clientset.CoreV1().Services(namespace).Delete(serviceName, nil)
	}
	return err
}

func isIP(input string) bool {
	pattern := "^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
	regEx := regexp.MustCompile(pattern)
	return regEx.MatchString(input)
}

func isHostName(input string) bool {
	pattern := "^(([a-zA-Z]|[a-zA-Z][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z]|[A-Za-z][A-Za-z0-9\\-]*[A-Za-z0-9])$"
	regEx := regexp.MustCompile(pattern)
	return regEx.MatchString(input)
}

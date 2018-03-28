package broker

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	v1alpha1 "k8s.io/api/settings/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func createBindingParameters(serviceInstance dbServiceInstance, bindingParameters map[string]interface{}) map[string]interface{} {

	ds := newDataSource(serviceInstance, bindingParameters)
	appType := i2s(bindingParameters["application-type"])

	switch appType {
	case "spring-boot": return ds.springboot()
	case "wildfly-swarm": return ds.wildflyswarm()
	case "nodejs": return ds.nodejs()
	case "other": return ds.other()
	}
	return ds.other()
}

func createPodPreset(serviceInstance dbServiceInstance, props map[string]interface{}) (bool, error){
	config, err := rest.InClusterConfig()
	if err != nil {
		return false, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}
	sourceName := i2s(serviceInstance.Parameters["source-name"])
	namespace := i2s(serviceInstance.Parameters["namespace"])

	labels :=  map[string]string{}
	labels["role"] = sourceName

	pp := v1alpha1.PodPreset{}
	pp.ObjectMeta.Name = sourceName
	pp.Spec.Selector = metav1.LabelSelector{MatchLabels: labels}
	//pp.Spec.Env =

	//
	//v1beta1.
	//scClient, err := v1beta1.NewForConfig(config)
	//scClient.ServiceBindings(namespace).
	//
	//clientset.serv
	//clientset.SettingsV1alpha1().PodPresets(namespace).Create()
}
package broker

import (
	osb "github.com/pmorie/go-open-service-broker-client/v2"
	"k8s.io/api/core/v1"
	"k8s.io/api/settings/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"errors"
	"fmt"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
	"time"
)

func createBindingParameters(serviceInstance dbServiceInstance,
	bindingParameters map[string]interface{}) map[string]interface{} {

	ds := newDataSource(serviceInstance, bindingParameters)
	appType := i2s(bindingParameters["application-type"])

	switch appType {
	case "spring-boot":
		return ds.springboot()
	case "wildfly-swarm":
		return ds.wildflyswarm()
	case "nodejs":
		return ds.nodejs()
	case "other":
		return ds.other()
	}
	return ds.other()
}

func podPresetAction(serviceInstance dbServiceInstance, name string, opKey osb.OperationKey,
	fn func(serviceInstance dbServiceInstance, appType string) error, after func(osb.OperationKey, error)) {
	// this is hack, really should be driven by some kind of event
	// I have found no such unless we can watch the service bindings
	// come through.
	var err error
	execute := fn(serviceInstance, name)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for t := range ticker.C {
			err = execute
			if err == nil {
				t.String()
				err = nil
				ticker.Stop()
				break
			}
		}
	}()
	time.Sleep(30 * time.Second)
	ticker.Stop()
	after(opKey, err)
}

func buildPodPreset(serviceInstance dbServiceInstance, podPresetName string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	sourceName := i2s(serviceInstance.Parameters["source-name"])
	namespace := i2s(serviceInstance.Parameters["namespace"])

	// lookthrough all the secrets and find the secret with source name property
	typeMetadata := metav1.TypeMeta{Kind: "Secret", APIVersion: "v1"}
	secretList, err := clientset.CoreV1().Secrets(namespace).List(metav1.ListOptions{TypeMeta: typeMetadata})
	if err != nil {
		fmt.Println("Failed to find the secret with source-name field" + i2s(err))
		return err
	}
	presetCreated := false
	for i := range secretList.Items {
		sn, ok := secretList.Items[i].Data["source-name"]
		if !ok {
			continue
		}
		srcName := string(sn)
		if srcName == sourceName {
			properties := secretList.Items[i].Data

			fmt.Println("Properties in selected secret" + i2s(properties))

			// create PodPreset
			labels := map[string]string{}
			labels["role"] = sourceName

			pp := v1alpha1.PodPreset{}
			pp.ObjectMeta.Name = podPresetName
			pp.Spec.Selector = metav1.LabelSelector{MatchLabels: labels}

			envs := []v1.EnvVar{}
			for key := range properties {
				if key == "source-name" {
					continue
				}

				ref := v1.LocalObjectReference{Name: secretList.Items[i].ObjectMeta.Name}
				secret := v1.SecretKeySelector{Key: key, LocalObjectReference: ref}
				value := v1.EnvVarSource{SecretKeyRef: &secret}
				key = strings.Replace(key, ".", "_", -1)
				key = strings.Replace(key, "-", "", -1)
				key = strings.ToUpper(key)
				env := v1.EnvVar{Name: key, ValueFrom: &value}
				envs = append(envs, env)
			}
			pp.Spec.Env = envs

			created, err := clientset.SettingsV1alpha1().PodPresets(namespace).Create(&pp)
			fmt.Println("created PodPreset" + i2s(pp))
			if err != nil {
				fmt.Println("Failed to create podpreset after create call " + i2s(err))
				return err
			}
			glog.Infof("created PodPreset" + i2s(created))
			fmt.Println("created PodPreset" + i2s(created))
			presetCreated = true
		}
	}
	if presetCreated {
		return nil
	}
	return errors.New("PodPreset is not created")
}

func removePodPreset(serviceInstance dbServiceInstance, podPresetName string) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	namespace := i2s(serviceInstance.Parameters["namespace"])
	typeMeta := metav1.TypeMeta{Kind: "PodPreset", APIVersion: "v1"}

	_, err = clientset.SettingsV1alpha1().PodPresets(namespace).Get(podPresetName, metav1.GetOptions{TypeMeta: typeMeta})
	if err != nil {
		return errors.New("failed to find the PodPreset with name " + podPresetName + " for removal.")
	}
	err = clientset.SettingsV1alpha1().PodPresets(namespace).Delete(podPresetName, &metav1.DeleteOptions{TypeMeta: typeMeta})
	return err
}

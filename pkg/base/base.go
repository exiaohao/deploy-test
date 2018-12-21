package base

import (
	"io/ioutil"
	"github.com/golang/glog"
	api_v1 "k8s.io/api/core/v1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes/scheme"
)

func InitializeKubeClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

// CreateNamespace for test
func CreateNamespace(kubeClient *kubernetes.Clientset, name string, showDetail bool) error {
	if name == "" {
		name = "test-namespace"
	}
	namespaceSpec := &api_v1.Namespace{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: name,
		},
	}
	_, err := kubeClient.CoreV1().Namespaces().Create(namespaceSpec)
	if showDetail && err == nil {
		glog.Info(CreateNamespaceSucceed(name))
	}
	return err
}

// RemoveNamespace after test
func RemoveNamespace(kubeClient *kubernetes.Clientset, name string, showDetail bool) error {
	if name == "" {
		name = "test-namespace"
	}
	err := kubeClient.CoreV1().Namespaces().Delete(name, &meta_v1.DeleteOptions{})
	if showDetail && err == nil {
		glog.Info(DestoryNamespace(name))
	}
	return err
}

// CreateDeployment
func CreateDeployment(kubeClient *kubernetes.Clientset, namespace string, fileFath string, showDetail bool) error {
	deoploymentFile, err := ioutil.ReadFile(fileFath)
	if err != nil {
		return CreateFailed("Read deployment file", err)
	}
	deploymentObj, _, err := scheme.Codecs.UniversalDeserializer().Decode(deoploymentFile, nil, nil)
	if err != nil {
		return CreateFailed("Parse deployment failed", err)
	}
	deployment := deploymentObj.(*extensions_v1beta1.Deployment)

	if namespace != "" {
		deployment.Namespace = namespace
	}

	_, err = kubeClient.ExtensionsV1beta1().Deployments(namespace).Create(deployment)
	if err != nil {
		return CreateFailed("Create deployment failed", err)
	}
	if showDetail {
		glog.Info(CreateDeploymentSucceed(deployment.Name))
	}
	return nil
}

// CreateService
func CreateService(kubeClient *kubernetes.Clientset, namespace string, fileFath string, showDetail bool) error {
	serviceFile, err := ioutil.ReadFile(fileFath)
	if err != nil {
		return CreateFailed("Read service file", err)
	}
	serviceObj, _, err := scheme.Codecs.UniversalDeserializer().Decode(serviceFile, nil, nil)
	if err != nil {
		return CreateFailed("Parse service failed", err)
	}
	service := serviceObj.(*api_v1.Service)

	if namespace != "" {
		service.Namespace = namespace
	}

	if _, err = kubeClient.CoreV1().Services(namespace).Create(service); err != nil {
		return CreateFailed("Create service failed", err)
	}

	if showDetail {
		glog.Info(CreateServiceSucceed(service.Name))
	}
	return nil
}
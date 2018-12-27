package base

import (
	"fmt"
	"github.com/golang/glog"
	"net/http"

	"io/ioutil"
	api_v1 "k8s.io/api/core/v1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
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
func CreateService(kubeClient *kubernetes.Clientset, namespace string, fileFath string, showDetail bool) (*api_v1.Service, error) {
	serviceFile, err := ioutil.ReadFile(fileFath)
	if err != nil {
		return nil, CreateFailed("Read service file", err)
	}
	serviceObj, _, err := scheme.Codecs.UniversalDeserializer().Decode(serviceFile, nil, nil)
	if err != nil {
		return nil, CreateFailed("Parse service failed", err)
	}
	service := serviceObj.(*api_v1.Service)

	if namespace != "" {
		service.Namespace = namespace
	}

	if _, err = kubeClient.CoreV1().Services(namespace).Create(service); err != nil {
		return nil, CreateFailed("Create service failed", err)
	}

	if showDetail {
		glog.Info(CreateServiceSucceed(service.Name))
	}
	return service, nil
}

// CheckServiceWorks
func CheckServiceWorks(kubeClient *kubernetes.Clientset, serviceObj *api_v1.Service, path string) error {
	service, err := kubeClient.CoreV1().Services(serviceObj.Namespace).Get(serviceObj.Name, meta_v1.GetOptions{})
	if err != nil {
		return BadServiceStatus(serviceObj.Name, fmt.Errorf("Get service %s failed", serviceObj.Name))
	}
	if path == "" {
		path = "/"
	}
	requestURL := fmt.Sprintf("http://%s:%d%s", service.Spec.ClusterIP, service.Spec.Ports[0].Port, path)
	resp, err := http.Get(requestURL)
	if err != nil {
		return BadServiceStatus(serviceObj.Name, err)
	}
	glog.Info(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return BadServiceStatus(serviceObj.Name, fmt.Errorf("Bad statusCode %d from %s", resp.StatusCode, requestURL))
	} else {
		glog.Infof("Request %s, returns statusCode: %d", requestURL, resp.StatusCode)
	}
	return nil
}

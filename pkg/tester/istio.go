package tester

import (
	"strings"

	"github.com/exiaohao/deploy-test/pkg/base"
	"github.com/golang/glog"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IstioTest struct {
	kubeClient    *kubernetes.Clientset
	namespace     string
	showDetail    bool
	testNamespace string
}

// Initialize istioTest
func (it *IstioTest) Initialize(opts InitOptions) {
	var err error

	if it.kubeClient, err = base.InitializeKubeClient(opts.KubeConfig); err != nil {
		glog.Fatalf("initialize kubeclient using %s: %v", opts.KubeConfig, err)
	}
	it.showDetail = true
	it.namespace = opts.Namespace
	it.testNamespace = "test-namespace"
}

// Run a istio test
func (it *IstioTest) Run() {
	if podsAvailableErr := it.checkPodsAvailable(); podsAvailableErr != nil{
		glog.Fatal(podsAvailableErr)
	}
	if projectErr := it.deploySimpleProject(); projectErr != nil {
		glog.Fatal(projectErr)
	}
}

// checkPodsAvailable check pods available: Running or Succeeded
func (it *IstioTest) checkPodsAvailable() error {
	pods, err := it.kubeClient.CoreV1().Pods(it.namespace).List(meta_v1.ListOptions{})
	if err != nil {
		glog.Fatalf("Get pods failed: %s", err)
	}
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case "Running":
			if it.showDetail {
				glog.Info(base.PodStatusOK(pod.Name, pod.Status.Phase))
			}
			continue
		case "Succeeded":
			if strings.Contains(pod.Name, "-post-install") || strings.Contains(pod.Name, "-cleanup-") {
				if it.showDetail {
					glog.Info(base.PodStatusOK(pod.Name, pod.Status.Phase))
				}
				continue
			}
			return base.BadPodStatus(pod.Name, pod.Status)
		default:
			return base.BadPodStatus(pod.Name, pod.Status)
		}
	}
	glog.Info(base.PodCheckPassed())
	return nil
}

func (it *IstioTest) deploySimpleProject() error {
	err := base.CreateNamespace(it.kubeClient, it.testNamespace, it.showDetail)
	if err != nil {
		return base.CreateNamespaceFailed(it.testNamespace, err)
	}
	defer base.RemoveNamespace(it.kubeClient, it.testNamespace, it.showDetail)

	if err := base.CreateDeployment(it.kubeClient, it.testNamespace,"./test-data/simple-project/deployment.yaml", it.showDetail); err != nil {
		return err
	}
	if err := base.CreateService(it.kubeClient, it.testNamespace, "./test-data/simple-project/service.yaml", it.showDetail); err != nil {
		return err
	}


	return nil
}


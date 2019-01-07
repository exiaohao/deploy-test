package tester

import (
	"fmt"
	"os"
	"strings"

	"github.com/exiaohao/deploy-test/pkg/base"
	"github.com/golang/glog"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type IstioTest struct {
	kubeClient     *kubernetes.Clientset
	namespace      string
	showDetail     bool
	runFullTest    bool
	displayErrFunc func(args ...interface{})
	testNamespace  string
}

// Initialize istioTest
func (it *IstioTest) Initialize(opts InitOptions) {
	var err error

	if it.kubeClient, err = base.InitializeKubeClient(opts.KubeConfig); err != nil {
		glog.Fatalf("initialize kubeclient using %s: %v", opts.KubeConfig, err)
	}

	it.namespace = opts.Namespace

	if os.Getenv("SHOW_DETAIL") == "TRUE" {
		it.showDetail = true
	} else {
		it.showDetail = false
	}

	if os.Getenv("IGNORE_FAILED") == "TRUE" {
		it.runFullTest = true
		it.displayErrFunc = glog.Info
	} else {
		it.runFullTest = false
		it.displayErrFunc = glog.Fatal
	}

	if os.Getenv("TEST_NAMESPACE") == "" {
		it.testNamespace = "test-namespace"
	} else {
		it.testNamespace = os.Getenv("TEST_NAMESPACE")
	}
}

// Run a istio test
func (it *IstioTest) Run() {
	// if podsAvailableErr := it.checkPodsAvailable(); podsAvailableErr != nil {
	// 	it.displayErrFunc(podsAvailableErr)
	// }
	// if projectErr := it.deploySimpleProject(); projectErr != nil {
	// 	it.displayErrFunc(projectErr)
	// }
	if bookinfoErr := it.bookinfo(); bookinfoErr != nil {
		it.displayErrFunc(bookinfoErr)
	}
}

// checkPodsAvailable check pods available: Running or Succeeded
func (it *IstioTest) checkPodsAvailable() error {
	pods, err := it.kubeClient.CoreV1().Pods(it.namespace).List(meta_v1.ListOptions{})
	if err != nil {
		// TODO RETURN FAILED
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
			if strings.Contains(pod.Name, "-post-install") || strings.Contains(pod.Name, "-cleanup-") || strings.Contains(pod.Name, "deploy-test") {
				if it.showDetail {
					glog.Info(base.PodStatusOK(pod.Name, pod.Status.Phase))
				}
				continue
			}
			return base.BadPodStatus(pod.Name, pod.Status)
		default:
			if strings.Contains(pod.Name, "deploy-test") {
				if it.showDetail {
					glog.Info(base.PodStatusWarn(pod))
				}
				continue
			}
			return base.BadPodStatus(pod.Name, pod.Status)
		}
	}
	glog.Info(base.CheckPassed("Pod check"))
	return nil
}

func (it *IstioTest) deploySimpleProject() error {
	if err := base.CreateNamespace(it.kubeClient, it.testNamespace, it.showDetail); err != nil {
		return base.CreateNamespaceFailed(it.testNamespace, err)
	}
	defer base.RemoveNamespace(it.kubeClient, it.testNamespace, it.showDetail)

	// TODO fix get path
	if err := base.CreateDeployment(it.kubeClient, it.testNamespace,
		"./test-data/simple-project/deployment.yaml", it.showDetail); err != nil {
		return err
	}
	// TODO fix get path
	if service, err := base.CreateService(it.kubeClient, it.testNamespace,
		"./test-data/simple-project/service.yaml", it.showDetail); err != nil {
		return err
	} else {
		if err := base.CheckServiceWorks(it.kubeClient, service, "/status/200", 1, 1, 3); err != nil {
			return err
		}
	}
	glog.Info(base.CheckPassed("Deploy simple project"))
	return nil
}

// bookinfo test
func (it *IstioTest) bookinfo() error {
	if err := base.CreateNamespace(it.kubeClient, it.testNamespace, it.showDetail); err != nil {
		return base.CreateNamespaceFailed(it.testNamespace, err)
	}
	defer base.RemoveNamespace(it.kubeClient, it.testNamespace, it.showDetail)

	deployments := []string{
		"deployment-detail-v1.yaml",
		"deployment-productpage-v1.yaml",
		"deployment-rating-v1.yaml",
		"deployment-reviews-v1.yaml",
		"deployment-reviews-v2.yaml",
		"deployment-reviews-v3.yaml",
	}
	services := []string{
		"service-detail.yaml",
		"service-productpage.yaml",
		"service-rating.yaml",
		"service-reviews.yaml",
	}

	for _, deployment := range deployments {
		file := fmt.Sprintf("./test-data/bookinfo/%s", deployment)
		if err := base.CreateDeployment(it.kubeClient, it.testNamespace, file, it.showDetail); err != nil {
			return err
		}
	}

	for _, service := range services {
		file := fmt.Sprintf("./test-data/bookinfo/%s", service)
		if _, err := base.CreateService(it.kubeClient, it.testNamespace, file, it.showDetail); err != nil {
			return err
		}
	}

	if err := base.WaitNamespacePodsReady(it.kubeClient, it.testNamespace, 5, 2); err != nil {
		return err
	}

	return nil
}

package base

import (
	"fmt"
	core_v1 "k8s.io/api/core/v1"
)

const (
	Info = "ℹ️ "
	StatusOK = "🆗"
	Succeed = "✅"
	Warn = "⚠️"
	Error = "🔴"
)

type BaseInfo struct {
	Code int32
	Level string
	Message string
}

func (e *BaseInfo) Error() string {
	return fmt.Sprintf("[%d] %s %s", e.Code, e.Level, e.Message)
}

// PodCheckPassed
func PodCheckPassed() *BaseInfo {
	return &BaseInfo{
		Code: 2000,
		Level: Succeed,
		Message: "Pod status check passed!",
	}
}

// PodStatusOK
func PodStatusOK(podName string, podPhase core_v1.PodPhase) *BaseInfo {
	return &BaseInfo{
		Code:2001,
		Level:StatusOK,
		Message: fmt.Sprintf("Pod %s status %s", podName ,podPhase),
	}
}

// NamespaceStatusOK
func NamespaceStatusOK() *BaseInfo {
	return &BaseInfo{
		Code: 2011,
		Level: StatusOK,
		Message: "Namespace status check passed!",
	}
}

// CreateNamespaceSucceed
func CreateNamespaceSucceed(name string) *BaseInfo {
	return &BaseInfo{
		Code: 2012,
		Level: Succeed,
		Message: fmt.Sprintf("Create namespace %s succeed", name),
	}
}

// DeploymentStatusOK
func DeploymentStatusOK(deploymentName string, deploymentStatus string) *BaseInfo {
	return &BaseInfo{
		Code: 2021,
		Level:StatusOK,
		Message: fmt.Sprintf("Deployment % status %s", deploymentName, deploymentStatus),
	}
}

// CreateDeploymentSucceed
func CreateDeploymentSucceed(deploymentName string) *BaseInfo {
	return &BaseInfo{
		Code: 2022,
		Level: Succeed,
		Message: fmt.Sprintf("Create deployment %s succeed", deploymentName),
	}
}

// CreateServiceSucceed
func CreateServiceSucceed(serviceName string) *BaseInfo {
	return &BaseInfo{
		Code: 2032,
		Level: Succeed,
		Message: fmt.Sprintf("Create service %s succeed", serviceName),
	}
}

// DestoryNamespace
func DestoryNamespace(namespace string) *BaseInfo {
	return &BaseInfo{
		Code: 4011,
		Level: Info,
		Message: fmt.Sprintf("Namesapce %s destroyed", namespace),
	}
}

// BadPodStatus
func BadPodStatus(podName string, podStatus core_v1.PodStatus) *BaseInfo {
	return &BaseInfo{
		Code: 5001,
		Level: Error,
		Message: fmt.Sprintf("Pod %s in bad status: %s", podName, podStatus.Phase),
	}
}

// CreateNamespaceFailed
func CreateNamespaceFailed(namespaceName string, error error) *BaseInfo {
	return &BaseInfo{
		Code: 5011,
		Level: Error,
		Message: fmt.Sprintf("Create test namespace %s failed: %s", namespaceName, error),
	}
}

// BadServiceStatus
func BadServiceStatus(serviceName string, error error) *BaseInfo {
	return &BaseInfo{
		Code: 5031,
		Level: Error,
		Message: fmt.Sprintf("Service %s is not work: %s", serviceName, error),
	}
}

// CreateFailed
func CreateFailed(action string, error error) *BaseInfo {
	return &BaseInfo{
		Code: 5091,
		Level: Error,
		Message: fmt.Sprintf("%s failed: %s", action, error),
	}
}

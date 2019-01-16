package util

import (
	"fmt"
)

func kubeCommand(subCommand, namespace, yamlFileName string, kubeconfig string) string {
	if namespace == "" {
		return fmt.Sprintf("kubectl %s -f %s --kubeconfig=%s", subCommand, yamlFileName, kubeconfig)
	}
	return fmt.Sprintf("kubectl %s -n %s -f %s --kubeconfig=%s", subCommand, namespace, yamlFileName, kubeconfig)
}

// KubeApply kubectl apply from file
func KubeApply(namespace, yamlFileName string, kubeconfig string) (string, error) {
	return Shell(kubeCommand("apply", namespace, yamlFileName, kubeconfig))
}

// KubeApplyContents kubectl apply from contents
func KubeApplyContents(namespace, yamlContents string, kubeconfig string) error {
	return Shell(kubeCommand(""))
}

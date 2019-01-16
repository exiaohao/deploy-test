package util

import "fmt"

func istioCommand(subCommand, yamlFileName, kubeConfig string) string {
	return fmt.Sprintf("istioctl %s -f %s --kubeconfig=%s", subCommand, yamlFileName, kubeConfig)
}

// IstioKubeInject istioctl kube-inject file
func IstioKubeInject(yamlFileName, kubeConfig string) (string, error) {
	return Shell(istioCommand("kube-inject", yamlFileName, kubeConfig))
}

package client

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// NewMetricsClientSet 创建一个新的 versioned.Clientset 实例
func NewMetricsClientSet() (*versioned.Clientset, error) {
	kubeconfig := "etc/ai-stage.yaml"
	//kubeconfig := "etc/ai-dx-test.yaml"
	//kubeconfig := "etc/npu-910b-test.yaml"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return versioned.NewForConfig(config)
}

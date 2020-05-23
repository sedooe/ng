package kubeclient

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	Client kubernetes.Interface
}

func NewFakeKubernetesClient(objects ...runtime.Object) KubernetesClient {
	return KubernetesClient{kubernetes.Interface(fake.NewSimpleClientset(objects...))}
}

func NewKubeClient() KubernetesClient {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return KubernetesClient{clientset}
}

package kubeclient

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func (kc KubernetesClient) GetOperationNamespace(fl *genericclioptions.ConfigFlags) (*v1.Namespace, error) {
	namespace := *fl.Namespace
	if namespace == "" {
		n, _, err := fl.ToRawKubeConfigLoader().Namespace()
		if err != nil {
			return nil, err
		}
		namespace = n
	}

	cli := kc.Client

	ns, err := cli.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func (kc KubernetesClient) UpdateNamespace(ns *v1.Namespace) error {
	_, err := kc.Client.CoreV1().Namespaces().Update(ns)
	if err != nil {
		return err
	}
	return nil
}

func (kc KubernetesClient) GetNamespaces(labelSelector string) ([]v1.Namespace, error) {
	namespaceList, err := kc.Client.CoreV1().Namespaces().List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}

	return namespaceList.Items, nil
}

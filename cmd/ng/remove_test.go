package main

import (
	"github.com/sedooe/ng/internal/kubeclient"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestRemoveNgForErrors(t *testing.T) {
	nsWithoutLabel := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}

	nsWithLabel := nsWithoutLabel
	nsWithLabel.Labels = map[string]string{label: testNg}

	tests := []struct {
		name string
		ns   *v1.Namespace
		o    *removeOptions
	}{
		{
			name: "namespace does not exist",
			ns:   &nsWithLabel,
			o: &removeOptions{
				cli: kubeclient.NewFakeKubernetesClient(),
			},
		},
		{
			name: "namespace does not belong to any group",
			ns:   &nsWithoutLabel,
			o: &removeOptions{
				cli: kubeclient.NewFakeKubernetesClient(&nsWithoutLabel),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.o.removeNg(tt.ns)
			if err == nil {
				t.Errorf("expected error didn't occur")
			}
		})
	}
}

func TestRemoveNgForSuccess(t *testing.T) {
	nsWithLabel := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "ns",
			Labels: map[string]string{label: testNg},
		},
	}

	nsWithMoreLabel := nsWithLabel
	nsWithMoreLabel.Labels = map[string]string{label: testNg, "a": "b"}

	tests := []struct {
		name string
		ns   *v1.Namespace
		o    *removeOptions
		want map[string]string
	}{
		{
			name: "namespace with ng label",
			ns:   &nsWithLabel,
			o: &removeOptions{
				cli: kubeclient.NewFakeKubernetesClient(&nsWithLabel),
			},
			want: map[string]string{},
		},
		{
			name: "namespace with ng and 1 more label",
			ns:   &nsWithMoreLabel,
			o: &removeOptions{
				cli: kubeclient.NewFakeKubernetesClient(&nsWithMoreLabel),
			},
			want: map[string]string{"a": "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.o.removeNg(tt.ns)
			if err != nil {
				t.Errorf("unexpected error occurred")
			}

			got := tt.ns.Labels

			eq := reflect.DeepEqual(tt.want, got)
			if !eq {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

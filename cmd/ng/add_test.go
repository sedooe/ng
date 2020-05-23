package main

import (
	"github.com/sedooe/ng/internal/kubeclient"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

const testNg = "backend"

func TestAddNgForErrors(t *testing.T) {
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
		o    *addOptions
	}{
		{
			name: "namespace does not exist",
			ns:   &nsWithoutLabel,
			o: &addOptions{
				name: testNg,
				cli:  kubeclient.NewFakeKubernetesClient(),
			},
		},
		{
			name: "namespace already belongs to a group",
			ns:   &nsWithLabel,
			o: &addOptions{
				name: testNg,
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithLabel),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.addNg(tt.ns)
			if err == nil {
				t.Errorf("expected error didn't occur")
			}
		})
	}
}

func TestAddNgForSuccess(t *testing.T) {
	nsWithoutLabel := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}

	nsWithLabel := nsWithoutLabel
	nsWithLabel.Labels = map[string]string{"a": "b"}

	tests := []struct {
		name string
		ns   *v1.Namespace
		o    *addOptions
		want map[string]string
	}{
		{
			name: "namespace with no label",
			ns:   &nsWithoutLabel,
			o: &addOptions{
				name: testNg,
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel),
			},
			want: map[string]string{label: testNg},
		},
		{
			name: "namespace with 1 other label",
			ns:   &nsWithLabel,
			o: &addOptions{
				name: testNg,
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithLabel),
			},
			want: map[string]string{"a": "b", label: testNg},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.addNg(tt.ns)
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

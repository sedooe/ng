package main

import (
	"github.com/sedooe/ng/internal/kubeclient"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"testing"
)

func TestShow(t *testing.T) {
	nsWithoutLabel := v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nsWithoutLabel",
		},
	}

	nsWithLabel := nsWithoutLabel
	nsWithLabel.Name = "nsWithLabel"
	nsWithLabel.Labels = map[string]string{"a": "b"}

	nsWithNgLabel := nsWithoutLabel
	nsWithNgLabel.Name = "nsWithNgLabel"
	nsWithNgLabel.Labels = map[string]string{"a": "b", label: testNg}

	nsWithNgLabel2 := nsWithNgLabel
	nsWithNgLabel2.Name = "nsWithNgLabel2"
	nsWithNgLabel2.Labels = map[string]string{"c": "d", label: "frontend"}

	nsWithNgLabel3 := nsWithNgLabel
	nsWithNgLabel3.Name = "nsWithNgLabel3"
	nsWithNgLabel3.Labels = map[string]string{label: testNg}

	tests := []struct {
		name string
		o    *showOptions
		want map[string][]string
	}{
		{
			name: "no ng with no command",
			o: &showOptions{
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithLabel),
			},
			want: make(map[string][]string),
		},
		{
			name: "no ng with a command when there is no ng namespace",
			o: &showOptions{
				name: testNg,
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithLabel),
			},
			want: make(map[string][]string),
		},
		{
			name: "no ng with a command when there are ng namespaces",
			o: &showOptions{
				name: "development",
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithNgLabel, &nsWithNgLabel2),
			},
			want: make(map[string][]string),
		},
		{
			name: "1 ng with a command",
			o: &showOptions{
				name: "frontend",
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithNgLabel, &nsWithNgLabel2),
			},
			want: map[string][]string{"frontend": {nsWithNgLabel2.Name}},
		},
		{
			name: "2 ng with no command",
			o: &showOptions{
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithNgLabel, &nsWithNgLabel3),
			},
			want: map[string][]string{testNg: {nsWithNgLabel.Name, nsWithNgLabel3.Name}},
		},
		{
			name: "2 ng with a command",
			o: &showOptions{
				name: testNg,
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithNgLabel, &nsWithNgLabel2, &nsWithNgLabel3),
			},
			want: map[string][]string{testNg: {nsWithNgLabel.Name, nsWithNgLabel3.Name}},
		},
		{
			name: "all ngs",
			o: &showOptions{
				cli:  kubeclient.NewFakeKubernetesClient(&nsWithoutLabel, &nsWithNgLabel, &nsWithNgLabel2, &nsWithNgLabel3),
			},
			want: map[string][]string{testNg: {nsWithNgLabel.Name, nsWithNgLabel3.Name}, "frontend": {nsWithNgLabel2.Name}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.o.show()
			if err != nil {
				t.Errorf("unexpected error occurred")
			}

			eq := reflect.DeepEqual(tt.want, got)
			if !eq {
				t.Errorf("want: %v, got: %v", tt.want, got)
			}
		})
	}
}

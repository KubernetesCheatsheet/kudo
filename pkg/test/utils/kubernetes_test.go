package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestNamespaced(t *testing.T) {
	fake := FakeDiscoveryClient()

	for _, test := range []struct {
		testName             string
		resource             runtime.Object
		namespaceShouldBeSet bool
		shouldError          bool
	}{
		{"namespaced resource", NewPod("hello", ""), true, false},
		{"namespace already set", NewPod("hello", "other"), true, false},
		{"not-namespaced resource", NewResource("v1", "Namespace", "hello", ""), false, false},
		{"non-existent resource", NewResource("v1", "Blah", "hello", ""), false, true},
	} {
		t.Run(test.testName, func(t *testing.T) {
			namespace := "world"

			m, _ := meta.Accessor(test.resource)

			actualName, actualNamespace, err := Namespaced(fake, test.resource, namespace)

			if test.shouldError {
				assert.NotNil(t, err)
				assert.Equal(t, "", actualName)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, m.GetName(), actualName)
			}

			if !test.namespaceShouldBeSet {
				namespace = ""
			}

			assert.Equal(t, namespace, actualNamespace)
			assert.Equal(t, namespace, m.GetNamespace())
		})
	}
}
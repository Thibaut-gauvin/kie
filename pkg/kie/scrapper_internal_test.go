package kie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//nolint:lll
func Test_parseImageName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		container corev1.Container
		pod       corev1.Pod
		expected  RunningImage
	}{
		{
			name:      "pod name ok",
			container: corev1.Container{},
			pod:       corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "ninja-pod-abcd"}},
			expected: RunningImage{
				name:   "",
				tag:    "",
				digest: "",
				pod:    "ninja-pod-abcd",
			},
		},
		{
			name:      "image: match only name",
			container: corev1.Container{Image: "alpine"},
			pod:       corev1.Pod{},
			expected: RunningImage{
				name:   "alpine",
				tag:    "",
				digest: "",
				pod:    "",
			},
		},
		{
			name:      "image: match name with tag",
			container: corev1.Container{Image: "alpine:3.19.1"},
			pod:       corev1.Pod{},
			expected: RunningImage{
				name:   "alpine",
				tag:    "3.19.1",
				digest: "",
				pod:    "",
			},
		},
		{
			name:      "image: match name with sha256",
			container: corev1.Container{Image: "alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b"},
			pod:       corev1.Pod{},
			expected: RunningImage{
				name:   "alpine",
				tag:    "",
				digest: "sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b",
				pod:    "",
			},
		},
		{
			name:      "image: match name with tag and sha256",
			container: corev1.Container{Image: "alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b"},
			pod:       corev1.Pod{},
			expected: RunningImage{
				name:   "alpine",
				tag:    "3.19.1",
				digest: "sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b",
				pod:    "",
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := parseImageName(test.container, test.pod)
			assert.Equal(t, test.expected, actual)
		})
	}
}

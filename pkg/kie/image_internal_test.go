package kie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

//nolint:lll
func Test_parseImageName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		container corev1.Container
		expected  RunningImage
	}{
		{
			name:      "empty match",
			container: corev1.Container{},
			expected: RunningImage{
				Name:               "",
				FullyQualifiedName: "",
				Tag:                "",
				Digest:             "",
				count:              1,
			},
		},
		{
			name:      "image: match only name",
			container: corev1.Container{Image: "alpine"},
			expected: RunningImage{
				Name:               "alpine",
				FullyQualifiedName: "alpine",
				Tag:                "",
				Digest:             "",
				count:              1,
			},
		},
		{
			name:      "image: match name with tag",
			container: corev1.Container{Image: "alpine:3.19.1"},
			expected: RunningImage{
				Name:               "alpine",
				FullyQualifiedName: "alpine:3.19.1",
				Tag:                "3.19.1",
				Digest:             "",
				count:              1,
			},
		},
		{
			name:      "image: match name with sha256",
			container: corev1.Container{Image: "alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b"},
			expected: RunningImage{
				Name:               "alpine",
				FullyQualifiedName: "alpine@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b",
				Tag:                "",
				Digest:             "sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b",
				count:              1,
			},
		},
		{
			name:      "image: match name, tag and sha256",
			container: corev1.Container{Image: "alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b"},
			expected: RunningImage{
				Name:               "alpine",
				FullyQualifiedName: "alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b",
				Tag:                "3.19.1",
				Digest:             "sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b",
				count:              1,
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := parseContainerImage(test.container)
			assert.Equal(t, test.expected, actual)
		})
	}
}

func Test_getImageFQN(t *testing.T) {
	t.Parallel()
	type input struct {
		imageName   string
		imageTag    string
		imageDigest string
	}
	tests := []struct {
		name     string
		input    input
		expected string
	}{
		{
			name:     "name only",
			input:    input{imageName: "nginx", imageTag: "", imageDigest: ""},
			expected: "nginx",
		},
		{
			name:     "name with tag",
			input:    input{imageName: "nginx", imageTag: "1.25.4-alpine", imageDigest: ""},
			expected: "nginx:1.25.4-alpine",
		},
		{
			name: "name with digest",
			input: input{
				imageName:   "nginx",
				imageTag:    "",
				imageDigest: "sha256:8ff7a5d76dd5a14c4d22d91ea45940bddc6384cdfe0b0348cfecdb9450c3143a",
			},
			expected: "nginx@sha256:8ff7a5d76dd5a14c4d22d91ea45940bddc6384cdfe0b0348cfecdb9450c3143a",
		},
		{
			name: "name with misspelled digest",
			input: input{
				imageName:   "nginx",
				imageTag:    "",
				imageDigest: "8ff7a5d76dd5a14c4d22d91ea45940bddc6384cdfe0b0348cfecdb9450c3143a",
			},
			expected: "nginx@sha256:8ff7a5d76dd5a14c4d22d91ea45940bddc6384cdfe0b0348cfecdb9450c3143a",
		},
		{
			name: "name, tag & digest",
			input: input{
				imageName:   "nginx",
				imageTag:    "1.25.4-alpine",
				imageDigest: "8ff7a5d76dd5a14c4d22d91ea45940bddc6384cdfe0b0348cfecdb9450c3143a",
			},
			expected: "nginx:1.25.4-alpine@sha256:8ff7a5d76dd5a14c4d22d91ea45940bddc6384cdfe0b0348cfecdb9450c3143a",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual := getImageFQN(test.input.imageName, test.input.imageTag, test.input.imageDigest)
			assert.Equal(t, test.expected, actual)
		})
	}
}

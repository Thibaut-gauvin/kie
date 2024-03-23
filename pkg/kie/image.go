package kie

import (
	"fmt"
	"regexp"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

var rxFrom = regexp.MustCompile(`^(?P<ref>(?P<image>[^:@\s]+):?(?P<tag>[^\s@]+)?@?(?P<digest>sha256:.*)?)$`)

type RunningImage struct {
	Name               string
	FullyQualifiedName string
	Tag                string
	Digest             string
	count              int
}

func parseContainerImage(container corev1.Container) RunningImage {
	runningImage := RunningImage{
		Name:  container.Image,
		count: 1,
	}
	if !rxFrom.MatchString(container.Image) {
		return runningImage
	}

	match := rxFrom.FindStringSubmatch(container.Image)
	result := make(map[string]string)
	for i, name := range rxFrom.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	runningImage.Name = result["image"]
	runningImage.Tag = result["tag"]
	runningImage.Digest = result["digest"]
	runningImage.FullyQualifiedName = getImageFQN(runningImage.Name, runningImage.Tag, runningImage.Digest)

	return runningImage
}

func getImageFQN(name, tag, digest string) string {
	imageFQN := name
	if tag != "" {
		imageFQN = fmt.Sprintf("%s:%s", imageFQN, tag)
	}
	if digest != "" {
		if !strings.HasPrefix(digest, "sha256:") {
			digest = fmt.Sprintf("sha256:%s", digest)
		}
		imageFQN = fmt.Sprintf("%s@%s", imageFQN, digest)
	}
	return imageFQN
}

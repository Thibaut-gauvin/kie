package kie

import (
	"context"
	"regexp"

	internalHandlers "github.com/Thibaut-gauvin/kie/internal/handlers"
	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/Thibaut-gauvin/kie/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var rxFrom = regexp.MustCompile(`^(?P<ref>(?P<image>[^:@\s]+):?(?P<tag>[^\s@]+)?@?(?P<digest>sha256:.*)?)$`)

type RunningImage struct {
	name   string
	tag    string
	digest string
	pod    string
}
type RunningImages []RunningImage

// scrapClusterImages is the function called by cron scheduler that is responsible to refresh prometheus metrics.
func scrapClusterImages(k8sClient kubernetes.Clientset) error {
	logger.Infof("tasks [%s] updated", clusterImageScrapJobName)

	pods, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		internalHandlers.UpdateHealthy(false)
		return err
	}

	var images RunningImages
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			images = append(images, parseImageName(container, pod))
		}
		for _, container := range pod.Spec.InitContainers {
			images = append(images, parseImageName(container, pod))
		}
	}

	for _, image := range images {
		metrics.UpdateRunningImagesGauge(
			image.name,
			image.tag,
			image.digest,
			image.pod,
		)
		logger.Debugf("image: %s, tag: %s, digest: %s, pod: %s",
			image.name,
			image.tag,
			image.digest,
			image.pod,
		)
	}

	internalHandlers.UpdateHealthy(true)
	return nil
}

func parseImageName(container corev1.Container, pod corev1.Pod) RunningImage {
	runningImage := RunningImage{
		name: container.Image,
		pod:  pod.Name,
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

	runningImage.name = result["image"]
	runningImage.tag = result["tag"]
	runningImage.digest = result["digest"]

	return runningImage
}

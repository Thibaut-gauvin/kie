package kie

import (
	"context"

	internalHandlers "github.com/Thibaut-gauvin/kie/internal/handlers"
	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/Thibaut-gauvin/kie/pkg/metrics"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// scrapClusterImages is the function called by cron scheduler that is responsible to refresh prometheus metrics.
func scrapClusterImages(k8sClient kubernetes.Clientset) error {
	logger.Infof("tasks [%s] updated", clusterImageScrapJobName)

	pods, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		internalHandlers.UpdateHealthy(false)
		return err
	}

	var containers []corev1.Container
	for _, pod := range pods.Items {
		containers = append(containers, pod.Spec.InitContainers...)
		containers = append(containers, pod.Spec.Containers...)
	}

	images := make(map[string]RunningImage)
	for _, container := range containers {
		image := parseContainerImage(container)

		_, isPresent := images[image.FullyQualifiedName]
		if isPresent {
			localImage := images[image.FullyQualifiedName]
			localImage.count++
			images[image.FullyQualifiedName] = localImage
		} else {
			images[image.FullyQualifiedName] = image
		}
	}

	metrics.ResetRunningImagesGauge()
	for _, image := range images {
		logger.Debugf("image: %s -> %d", image.FullyQualifiedName, image.count)
		metrics.UpdateRunningImagesGauge(image.FullyQualifiedName, image.Tag, image.Digest, image.count)
	}

	internalHandlers.UpdateHealthy(true)
	return nil
}

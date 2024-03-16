package kie

import (
	"context"

	internalHandlers "github.com/Thibaut-gauvin/kie/internal/handlers"
	"github.com/Thibaut-gauvin/kie/internal/logger"
	"github.com/Thibaut-gauvin/kie/pkg/metrics"
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

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			metrics.ObserveClusterImagesCount(container.Image, container.Name, pod.Name, pod.Namespace)
			logger.Debugf("image: %s, pod: %s, namespace: %s", container.Image, pod.Name, pod.Namespace)
		}

		for _, container := range pod.Spec.InitContainers {
			metrics.ObserveClusterImagesCount(container.Image, container.Name, pod.Name, pod.Namespace)
			logger.Debugf("image: %s, pod: %s, namespace: %s", container.Image, pod.Name, pod.Namespace)
		}
	}

	internalHandlers.UpdateHealthy(true)
	return nil
}

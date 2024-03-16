package kie

import (
	"context"
	"time"

	"github.com/Thibaut-gauvin/kie/internal/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// scrapClusterImages is the function called by cron scheduler that is responsible to refresh prometheus metrics.
func scrapClusterImages(k8sClient kubernetes.Clientset) error {
	logger.Infof("[%s] %s", clusterImageScrapJobName, time.Now().Format(time.RFC3339))

	pods, err := k8sClient.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			logger.Debugf(
				"image: %s, pod: %s, namespace: %s",
				container.Image, pod.Name, pod.Namespace,
			)
		}

		for _, container := range pod.Spec.InitContainers {
			logger.Debugf(
				"image: %s, pod: %s, namespace: %s",
				container.Image, pod.Name, pod.Namespace,
			)
		}
	}

	return nil
}

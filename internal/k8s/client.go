package k8s

import (
	"context"
	"fmt"
	"path"

	"github.com/Thibaut-gauvin/kie/internal/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// NewClient creates a k8s client instance, deciding whether "In-Cluster" or "Out-of-Cluster"
// authentication method must be used.
func NewClient(kubeconfig string) (*kubernetes.Clientset, error) {
	// First, try to "in-cluster" method
	client, err := NewInCluster()
	if err != nil {
		logger.Debugf("%v", err)
	} else {
		// Then, try listing namespaces
		namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			logger.Debugf("%v", err)
		}
		if namespaces != nil {
			return client, nil
		}
	}

	// Finally try "out-of-cluster"
	client, err = NewOutOfCluster(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return client, nil
}

// NewInCluster creates a KubeCli instance authenticated to cluster using service-account.
func NewInCluster() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s 'in-cluster' client: %w", err)
	}

	return kubernetes.NewForConfig(config)
}

// NewOutOfCluster creates a KubeCli instance authenticated to cluster using kubeconfig file.
func NewOutOfCluster(kubeconfig string) (*kubernetes.Clientset, error) {
	// If kubeconfig path is empty, try to load them from default location
	if kubeconfig == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = path.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create k8s 'out-of-cluster' client: %w", err)
	}
	return kubernetes.NewForConfig(config)
}

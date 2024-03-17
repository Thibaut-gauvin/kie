package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var runningImagesGauge = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "kie",
		Subsystem: "cluster",
		Name:      "images_running",
		Help:      "k8s cluster running images",
	},
	[]string{"name", "tag", "digest", "pod"},
)

func UpdateRunningImagesGauge(name, tag, digest, pod string) {
	runningImagesGauge.With(
		prometheus.Labels{
			"name":   name,
			"tag":    tag,
			"digest": digest,
			"pod":    pod,
		},
	).Set(1)
}

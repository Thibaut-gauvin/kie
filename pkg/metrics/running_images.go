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
		Help:      "Gauge counting images running on cluster",
	},
	[]string{"name", "tag", "digest"},
)

func UpdateRunningImagesGauge(fullyQualifiedName, tag, digest string, count int) {
	runningImagesGauge.With(
		prometheus.Labels{
			"name":   fullyQualifiedName,
			"tag":    tag,
			"digest": digest,
		},
	).Set(float64(count))
}

func ResetRunningImagesGauge() {
	runningImagesGauge.Reset()
}

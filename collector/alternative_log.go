package collector

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"

	"strconv"
	"strings"
)

const (
	alterNativeLogSubsystem = "alternativelog"
)

func init() {
	registerCollector(alterNativeLogSubsystem, defaultEnabled, NewalterNativeCollector)
}

type alterNativeCollector struct {
	logger log.Logger
}

func NewalterNativeCollector(logger log.Logger) (Collector, error) {
	return &alterNativeCollector{logger}, nil
}

func (c *alterNativeCollector) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType
	metricType = prometheus.GaugeValue
	output := alterNativeLog()
	for _, line := range strings.Split(output, "\n") {
		l := strings.Split(line, ":")
		if len(l) != 2 {
			continue
		}
		name := strings.TrimSpace(l[0])
		value := strings.TrimSpace(l[1])
		v, _ := strconv.Atoi(value)
		name = strings.Replace(name, "-", "_", -1)
		level.Debug(c.logger).Log("msg", "Set errLog", "name", name, "value", value)

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, alterNativeLogSubsystem, name),
				fmt.Sprintf("handle /var/log/alternatives.log info err log %s.", name),
				nil, nil,
			),
			metricType, float64(v),
		)
	}
	return nil
}

func alterNativeLog() string {
	alterNativeLogCmd := ` awk '!seen[$0]++ {count++; print "info_demo" count ":", length}' /var/log/alternatives.log`
	cmd := exec.Command("sh", "-c", alterNativeLogCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(output)
}

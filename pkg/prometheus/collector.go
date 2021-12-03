package prometheus

import (
	"cloud_1400/pkg/cast"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"runtime"
)

type metricTypes int

const (
	Load = iota + 1
	Success
	Error
)

type Collector struct {
	apiMetrics    map[string]*APIMetric
	customMetrics map[string]*Metric
}

type APIMetric struct {
	apiMetric map[metricTypes]*Metric
	stats     *Stats
}

type Metric struct {
	Name        string
	Description string
	ValueType   prometheus.ValueType
	ConstLabel  prometheus.Labels
	pDesk       *prometheus.Desc
	GetValue    func() float64
}

func NewCollector() *Collector {
	collector := new(Collector)
	collector.apiMetrics = make(map[string]*APIMetric)
	collector.customMetrics = make(map[string]*Metric)
	collector.setDefaultMetric()
	return collector
}
func (collector *Collector) AddAPIMetric(metric Metric, stats *Stats) {
	if _, found := collector.apiMetrics[metric.Name]; found {
		zap.L().Debug("repetitious metric", zap.String("metric name", metric.Name))
		return
	}
	if stats == nil {
		zap.L().Error("stats is nil")
		return
	}
	if metric.ValueType == 0 {
		metric.ValueType = prometheus.CounterValue
	}
	collector.apiMetrics[metric.Name] = &APIMetric{
		apiMetric: map[metricTypes]*Metric{
			Load: {
				Name:        metric.Name + "_load_request",
				Description: "load stats of " + metric.Name + " api",
				ValueType:   metric.ValueType,
				ConstLabel:  metric.ConstLabel,
				pDesk:       prometheus.NewDesc(metric.Name+"_load_request", metric.Description, nil, metric.ConstLabel),
			},
			Success: {
				Name:        metric.Name + "_success_request",
				Description: "success stats of " + metric.Name + " api",
				ValueType:   metric.ValueType,
				ConstLabel:  metric.ConstLabel,
				pDesk:       prometheus.NewDesc(metric.Name+"_success_request", metric.Description, nil, metric.ConstLabel),
			},
			Error: {
				Name:        metric.Name + "_error_request",
				Description: "error stats of " + metric.Name + " api",
				ValueType:   metric.ValueType,
				ConstLabel:  metric.ConstLabel,
				pDesk:       prometheus.NewDesc(metric.Name+"_error_request", metric.Description, nil, metric.ConstLabel),
			},
		},
		stats: stats,
	}
}
func (collector *Collector) BulkAddAPIMetric(metrics []Metric, stats []*Stats) {
	if len(metrics) != len(stats) {
		panic("metrics and stats not match")
	}
	for i := 0; i < len(metrics); i++ {
		collector.AddAPIMetric(metrics[i], stats[i])
	}
}
func (collector *Collector) AddCustomMetric(metric Metric) {
	if _, found := collector.customMetrics[metric.Name]; found {
		zap.L().Debug("repetitious metric", zap.String("metric name", metric.Name))
		return
	}
	if metric.GetValue == nil {
		zap.L().Error("getValue is nil")
		return
	}
	if metric.ValueType == 0 {
		metric.ValueType = prometheus.CounterValue
	}
	collector.customMetrics[metric.Name] = &Metric{
		Name:        metric.Name,
		Description: metric.Description,
		ValueType:   metric.ValueType,
		ConstLabel:  metric.ConstLabel,
		pDesk:       prometheus.NewDesc(metric.Name, metric.Description, nil, metric.ConstLabel),
		GetValue:    metric.GetValue,
	}
}
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range collector.apiMetrics {
		copyMetric := *metric
		m := &copyMetric
		ch <- m.GetLoadMetric().pDesk
		ch <- m.GetSuccessMetric().pDesk
		ch <- m.GetErrorMetric().pDesk
	}
	for _, metric := range collector.customMetrics {
		copyMetric := *metric
		m := &copyMetric
		ch <- m.pDesk
	}
}
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, metric := range collector.apiMetrics {
		loadMetric := metric.GetLoadMetric()
		successMetric := metric.GetSuccessMetric()
		errMetric := metric.GetErrorMetric()
		ch <- prometheus.MustNewConstMetric(loadMetric.pDesk, loadMetric.ValueType, metric.stats.GetLoads())
		ch <- prometheus.MustNewConstMetric(successMetric.pDesk, successMetric.ValueType, metric.stats.GetSuccess())
		ch <- prometheus.MustNewConstMetric(errMetric.pDesk, errMetric.ValueType, metric.stats.GetError())
	}
	for _, metric := range collector.customMetrics {
		copyMetric := *metric
		m := &copyMetric
		ch <- prometheus.MustNewConstMetric(m.pDesk, m.ValueType, m.GetValue())
	}
}
func (collector *Collector) setDefaultMetric() {
	memoryMetric := Metric{
		Name:        "process_mem_alloc_bytes_total",
		Description: "Show memory usage",
		ValueType:   prometheus.GaugeValue,
		ConstLabel:  nil,
		GetValue: func() float64 {
			var memory runtime.MemStats
			runtime.ReadMemStats(&memory)
			t, _ := cast.ToFloat64(memory.Alloc)
			return t
		},
	}
	memoryMetric.pDesk = prometheus.NewDesc(memoryMetric.Name, memoryMetric.Description, nil, memoryMetric.ConstLabel)
	collector.AddCustomMetric(memoryMetric)
}
func (a APIMetric) GetLoadMetric() *Metric {
	return a.apiMetric[Load]
}
func (a APIMetric) GetSuccessMetric() *Metric {
	return a.apiMetric[Success]
}
func (a APIMetric) GetErrorMetric() *Metric {
	return a.apiMetric[Error]
}

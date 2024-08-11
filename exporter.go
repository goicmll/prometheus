package prometheus

import (
	"strconv"
	"strings"
)

const (
	// defaultMetricCount 预估的指标个数， 即预估的本次获取的所有样本的指标名去重的个数
	defaultMetricCount uint64 = 2
	// defaultSampleCount 预估每个指标名有多少个样本
	defaultSampleCount uint64 = 8
)

// Exporter 指标导出
type Exporter struct {
	// metricCount 预估的指标个数， 即预估的本次获取的所有样本的指标名去重的个数
	metricCount uint64
	// sampleCount 预估每个指标名有多少个样本
	sampleCount   uint64
	MetricSamples map[string][]*Sample
}

// NewExporter 创建并返回一个新的Exporter实例。
// 它根据提供的metric数量和总样本数量来初始化Exporter。
// 如果metricCount和totalSampleCount都非零，且totalSampleCount大于等于metricCount，
// 则使用这些值来计算每个metric的样本数量；否则，使用默认值。
// 返回的Exporter实例包含一个初始化的MetricSamples map，用于存储metric名称和样本的关联。
func NewExporter(metricCount, totalSampleCount uint64) *Exporter {
	// 初始化默认的每个metric的样本数量和metric数量
	var preMetricSampleCount = defaultSampleCount
	var metricCount_ = defaultMetricCount

	// 如果提供了有效的metricCount和totalSampleCount，则使用它们来计算每个metric的样本数量
	if metricCount != 0 && totalSampleCount != 0 && totalSampleCount >= metricCount {
		metricCount_ = metricCount
		preMetricSampleCount = totalSampleCount / metricCount
	}

	// 使用计算出的metric数量和样本数量来创建Exporter实例
	// 初始化MetricSamples map，准备存储metric名称和样本的关联
	return &Exporter{metricCount: metricCount_, sampleCount: preMetricSampleCount, MetricSamples: make(map[string][]*Sample, metricCount)}
}

// String 方法返回Exporter的字符串表示。
// 它遍历所有的MetricSamples，并构建一个包含所有指标数据的字符串。
// 该方法主要用于导出指标数据为字符串格式，以便于日志记录或调试目的。
func (e *Exporter) String() string {
	var builder strings.Builder
	for _, samples := range e.MetricSamples {
		if len(samples) == 0 {
			continue
		}
		// 写入help
		builder.WriteString("\n")
		builder.WriteString("# HELP ")
		builder.WriteString(samples[0].MetricName)
		builder.WriteString(" ")
		builder.WriteString(samples[0].Help)
		builder.WriteString("\n")

		// 写入 type
		builder.WriteString("# TYPE ")
		builder.WriteString(samples[0].MetricName)
		builder.WriteString(" ")
		builder.WriteString(samples[0].Type)
		builder.WriteString("\n")

		for _, sample := range samples {
			// 写入样本
			builder.WriteString(sample.MetricName)
			builder.WriteString(mapToStr(sample.Labels))
			builder.WriteString(" ")
			builder.WriteString(strconv.FormatFloat(sample.Value, 'f', sample.ValuePrecision, 64))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

// AddSamples 添加指标样本
func (e *Exporter) AddSamples(ss ...*Sample) {
	for _, s := range ss {
		if _, ok := e.MetricSamples[s.MetricName]; ok {
			// 已存在,追加
			e.MetricSamples[s.MetricName] = append(e.MetricSamples[s.MetricName], s)
		} else {
			// 不存在,先初始化,再追加
			e.MetricSamples[s.MetricName] = make([]*Sample, 0, e.sampleCount)
			e.MetricSamples[s.MetricName] = append(e.MetricSamples[s.MetricName], s)
		}
	}
}

// Merge 合并 exporter
func (e *Exporter) Merge(es ...*Exporter) {
	for _, ee := range es {
		for _, ms := range ee.MetricSamples {
			e.AddSamples(ms...)
		}
	}
}

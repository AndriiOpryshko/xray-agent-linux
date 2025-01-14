package collectors

import (
	"errors"

	"xray-agent-linux/conf"
	"xray-agent-linux/dto"
	"xray-agent-linux/logger"
)

type MemoryDataSource interface {
	GetData() (*dto.MemoryInfo, error)
}

type MemoryInfoCollector struct {
	Config     *conf.MemoryInfoConf
	DataSource MemoryDataSource
}

func NewMemoryInfoCollector(cfg *conf.CollectorsConf, dataSource MemoryDataSource) dto.Collector {
	if cfg == nil || dataSource == nil {
		logger.LogWarning(logger.CollectorInitPrefix, errors.New("memory info collector init params error"))
		return nil
	}

	// exit if collector disabled
	if cfg.MemoryInfo == nil || !cfg.MemoryInfo.Enabled {
		return nil
	}

	return &MemoryInfoCollector{
		Config:     cfg.MemoryInfo,
		DataSource: dataSource,
	}
}

func (c *MemoryInfoCollector) GetName() string {
	return dto.CollectorNameMemoryInfo
}

func (c *MemoryInfoCollector) Collect() ([]dto.Metric, error) {
	memoryInfo, err := c.DataSource.GetData()
	if err != nil {
		return nil, err
	}

	attrs := []dto.MetricAttribute{
		{
			Name:  dto.ResourceAttr,
			Value: dto.ResourceMemory,
		},
	}

	return []dto.Metric{
		{
			Name:       dto.MetricMemoryTotal,
			Value:      memoryInfo.MemTotal,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemoryFree,
			Value:      memoryInfo.MemFree,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemoryAvailable,
			Value:      memoryInfo.MemAvailable,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemoryUsed,
			Value:      memoryInfo.MemTotal - memoryInfo.MemAvailable,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemoryBuffers,
			Value:      memoryInfo.Buffers,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemoryCached,
			Value:      memoryInfo.Cached,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemorySwapTotal,
			Value:      memoryInfo.SwapTotal,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricMemorySwapFree,
			Value:      memoryInfo.SwapFree,
			Attributes: attrs,
		},
	}, nil
}

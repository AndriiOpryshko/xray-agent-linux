package collectors

import (
	"errors"
	"strconv"

	"xray-agent-linux/conf"
	"xray-agent-linux/dto"
	"xray-agent-linux/logger"
)

type CPUInfoDataSource interface {
	GetData() (*dto.CPUInfo, error)
}

type CPUInfoCollector struct {
	Config     *conf.CPUInfoConf
	DataSource CPUInfoDataSource
}

func NewCpuInfoCollector(cfg *conf.CollectorsConf, dataSource CPUInfoDataSource) dto.Collector {
	if cfg == nil || dataSource == nil {
		logger.LogWarning(logger.CollectorInitPrefix, errors.New("cpu info collector init params error"))

		return nil
	}

	// exit if collector disabled
	if cfg.CPUInfo == nil || !cfg.CPUInfo.Enabled {
		return nil
	}

	return &CPUInfoCollector{
		Config:     cfg.CPUInfo,
		DataSource: dataSource,
	}
}

func (c *CPUInfoCollector) GetName() string {
	return dto.CollectorNameCPUInfo
}

func (c *CPUInfoCollector) Collect() ([]dto.Metric, error) {
	cpuInfo, err := c.DataSource.GetData()
	if err != nil {
		return nil, err
	}

	metrics := make([]dto.Metric, 0, len(cpuInfo.CPU))

	for _, cpu := range cpuInfo.CPU {
		metrics = append(metrics, dto.Metric{
			Name:     dto.MetricCPUInfoMHz,
			Attributes: []dto.MetricAttribute{
				{
					Name:  dto.ResourceAttr,
					Value: dto.ResourceCPUInfo,
				},
				{
					Name:  dto.SetNameCPUInfoProcessor,
					Value: strconv.Itoa(cpu.ProcessorNumber),
				},
			},
			Value: cpu.MHz,
		})
	}

	return metrics, nil
}

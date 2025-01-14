package collectors

import (
	"xray-agent-linux/conf"
	"xray-agent-linux/dto"
)

type NginxStubStatusCollector struct {
	Config     *conf.NginxStubStatus
	DataSource NginxStubStatusDataSource
}

func NewNginxStubStatusCollector(cfg *conf.CollectorsConf, dataSource NginxStubStatusDataSource) dto.Collector {
	if cfg == nil || dataSource == nil {
		return nil
	}

	// exit if collector disabled
	if cfg.NginxStubStatus == nil || !cfg.NginxStubStatus.Enabled {
		return nil
	}

	return &NginxStubStatusCollector{
		Config:     cfg.NginxStubStatus,
		DataSource: dataSource,
	}
}

type NginxStubStatusDataSource interface {
	GetData() (*dto.NginxStubStatus, error)
}

func (c *NginxStubStatusCollector) GetName() string {
	return dto.CollectorNameNginx
}

func (c *NginxStubStatusCollector) Collect() ([]dto.Metric, error) {
	data, err := c.DataSource.GetData()
	if err != nil {
		return nil, err
	}

	attrs := []dto.MetricAttribute{
		{
			Name:  dto.ResourceAttr,
			Value: dto.ResourceNginxStubStatus,
		},
	}

	return []dto.Metric{
		{
			Name:       dto.MetricNginxStubStatusActive,
			Value:      data.Active,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricNginxStubStatusAccepts,
			Value:      data.Accepts,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricNginxStubStatusHandled,
			Value:      data.Handled,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricNginxStubStatusRequests,
			Value:      data.Requests,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricNginxStubStatusReading,
			Value:      data.Reading,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricNginxStubStatusWriting,
			Value:      data.Writing,
			Attributes: attrs,
		},
		{
			Name:       dto.MetricNginxStubStatusWaiting,
			Value:      data.Waiting,
			Attributes: attrs,
		},
	}, nil
}

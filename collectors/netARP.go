package collectors

import (
	"errors"
	"strings"

	"xray-agent-linux/conf"
	"xray-agent-linux/dto"
	"xray-agent-linux/logger"
)

type NetARPDataSource interface {
	GetData() ([]dto.ARPEntry, error)
}

type NetARPCollector struct {
	Config     *conf.NetARPConf
	DataSource NetARPDataSource
}

func NewNetARPCollector(cfg *conf.CollectorsConf, dataSource NetARPDataSource) dto.Collector {
	if cfg == nil || dataSource == nil {
		logger.LogWarning(logger.CollectorInitPrefix, errors.New("net arp collector init params error"))
		return nil
	}

	// exit if collector disabled
	if cfg.NetARP == nil || !cfg.NetARP.Enabled {
		return nil
	}

	return &NetARPCollector{
		Config:     cfg.NetARP,
		DataSource: dataSource,
	}
}

func (c *NetARPCollector) GetName() string {
	return dto.CollectorNameNetARP
}

func (c *NetARPCollector) Collect() ([]dto.Metric, error) {
	netArp, err := c.getNetArp()
	if err != nil {
		return nil, err
	}

	metrics := make([]dto.Metric, 0, len(netArp.Entries)+len(netArp.IncompleteEntries))

	for devName, value := range netArp.Entries {
		devName = strings.ReplaceAll(devName, ".", "_")

		metrics = append(metrics,
			dto.Metric{
				Name:     dto.MetricNetARPEntries,
				Value:    value,
				Attributes: []dto.MetricAttribute{
					{
						Name:  dto.ResourceAttr,
						Value: dto.ResourceNetARP,
					},
					{
						Name:  dto.SetNameNetARPInterface,
						Value: devName,
					},
				},
			},
		)
	}

	for devName, value := range netArp.IncompleteEntries {
		devName = strings.ReplaceAll(devName, ".", "_")

		metrics = append(metrics,
			dto.Metric{
				Name:     dto.MetricNetARPIncompleteEntries,
				Value:    value,
				Attributes: []dto.MetricAttribute{
					{
						Name:  dto.ResourceAttr,
						Value: dto.ResourceNetARP,
					},
					{
						Name:  dto.SetNameNetARPInterface,
						Value: devName,
					},
				},
			},
		)
	}

	return metrics, nil
}

func (c *NetARPCollector) getNetArp() (*dto.NetArp, error) {
	arpTable, err := c.DataSource.GetData()
	if err != nil {
		return nil, err
	}

	var out dto.NetArp
	out.Entries = make(map[string]uint)
	out.IncompleteEntries = make(map[string]uint)

	out.Entries["Total"] = 0
	out.IncompleteEntries["Total"] = 0

	for _, entry := range arpTable {
		if _, ok := out.Entries[entry.Device]; !ok {
			out.Entries[entry.Device] = 0
		}

		if _, ok := out.IncompleteEntries[entry.Device]; !ok {
			out.IncompleteEntries[entry.Device] = 0
		}

		out.Entries["Total"]++
		out.Entries[entry.Device]++

		if entry.HWAddress == "00:00:00:00:00:00" {
			out.IncompleteEntries["Total"]++
			out.IncompleteEntries[entry.Device]++
		}
	}

	return &out, err
}

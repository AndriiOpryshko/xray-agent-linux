package collectors

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"xray-agent-linux/conf"
	"xray-agent-linux/dto"
	"xray-agent-linux/logger"
)

type CmdDataSource interface {
	RunPipeLine(pipeLine []*exec.Cmd) (string, string, error)
}

type CmdCollector struct {
	Config     *conf.CMDConf
	DataSource CmdDataSource
}

func NewCmdCollector(cfg *conf.CollectorsConf, dataSource CmdDataSource) dto.Collector {
	if cfg == nil || dataSource == nil {
		logger.LogWarning(logger.CollectorInitPrefix, errors.New("cmd collector init params error"))

		return nil
	}

	// exit if collector disabled
	if cfg.CMD == nil || !cfg.CMD.Enabled {
		return nil
	}

	return &CmdCollector{
		Config:     cfg.CMD,
		DataSource: dataSource,
	}
}

func (c *CmdCollector) GetName() string {
	return dto.CollectorNameCMD
}

func (c *CmdCollector) Collect() ([]dto.Metric, error) {
	metrics := make([]dto.Metric, 0, len(c.Config.Metrics))

	for i, _ := range c.Config.Metrics {
		err := c.processPipeLine(&c.Config.Metrics[i], c.Config.Timeout, &metrics)
		if err != nil {
			logger.LogWarning(dto.CollectorNameCMD, err)
		}
	}

	return metrics, nil
}

func (c *CmdCollector) processPipeLine(cfg *conf.CMDMetricConf, timeout int, out *[]dto.Metric) error {
	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	pipeLine := make([]*exec.Cmd, 0, len(cfg.PipeLine))

	for _, cmd := range cfg.PipeLine {
		if len(cmd) == 1 {
			pipeLine = append(pipeLine, exec.CommandContext(ctx, cmd[0]))
		}

		if len(cmd) >= 2 {
			pipeLine = append(pipeLine, exec.CommandContext(ctx, cmd[0], cmd[1:]...))
		}
	}

	stdout, stderr, err := c.DataSource.RunPipeLine(pipeLine)
	if err != nil {
		return err
	}

	// Timeout
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("command timed out")
	}

	// Check stderr
	if stderr != "" {
		return fmt.Errorf("stderr is not empty: '%s'", stderr)
	}

	values := strings.Split(strings.TrimSpace(stdout), cfg.Delimiter)

	if len(values) != len(cfg.Names) {
		return fmt.Errorf("metric count mismatch: config=%v, output=%v", len(cfg.Names), len(values))
	}

	for i, name := range cfg.Names {
		// skip ignored values
		if name == "-" {
			continue
		}

		*out = append(*out, dto.Metric{
			Name: cfg.Names[i],
			Attributes: append([]dto.MetricAttribute{
				{
					Name:  dto.ResourceAttr,
					Value: dto.ResourceCMD,
				},
			}, cfg.Attributes...),
			Value: values[i],
		})
	}

	return nil
}

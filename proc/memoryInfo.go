package proc

import (
	"strconv"
	"strings"

	"xray-agent-linux/dto"
	"xray-agent-linux/logger"
	"xray-agent-linux/reader"
)

type memoryDataSource struct {
	filePath  string
	logPrefix string
}

func NewMemoryDataSource(filePath, logPrefix string) *memoryDataSource {
	if filePath == "" {
		return nil
	}

	return &memoryDataSource{
		filePath:  filePath,
		logPrefix: logPrefix,
	}
}

func (ds *memoryDataSource) GetData() (*dto.MemoryInfo, error) {
	memoryInfo := dto.MemoryInfo{}

	lines, err := reader.ReadMultilineFile(ds.filePath, ds.logPrefix)
	if err != nil {
		return nil, err
	}

	for _, v := range lines {
		fields := strings.Fields(v)
		// skip incorrect lines
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			memoryInfo.MemTotal, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "MemTotal", err)
				// mandatory field
				return nil, err
			}
		case "MemFree:":
			memoryInfo.MemFree, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "MemFree", err)
				// mandatory field
				return nil, err
			}
		case "MemAvailable:":
			memoryInfo.MemAvailable, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "MemAvailable", err)
			}
		case "Buffers:":
			memoryInfo.Buffers, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "Buffers", err)
			}
		case "Cached:":
			memoryInfo.Cached, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "Cached", err)
			}
		case "SwapTotal:":
			memoryInfo.SwapTotal, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "SwapTotal", err)
			}
		case "SwapFree:":
			memoryInfo.SwapFree, err = strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				logger.LogReadFileFieldError(ds.logPrefix, ds.filePath, "SwapFree", err)
			}
		}
	}

	// if MemAvailable not parsed
	// calculate it
	if memoryInfo.MemAvailable == 0 {
		memoryInfo.MemAvailable = memoryInfo.MemTotal - memoryInfo.MemFree + memoryInfo.Cached + memoryInfo.Buffers
	}

	return &memoryInfo, nil
}

package agent

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
)

type NvidiaSMI struct {
	Path string
}

func FindNvidiaSMI() (NvidiaSMI, bool) {
	path, err := exec.LookPath("nvidia-smi")
	return NvidiaSMI{Path: path}, err == nil
}

func (n NvidiaSMI) Load(ctx context.Context) (map[int]float64, error) {
	out, err := exec.CommandContext(ctx, n.Path, "--query-gpu=index,utilization.gpu", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return nil, err
	}
	return parseLoad(string(out)), nil
}

func parseLoad(out string) map[int]float64 {
	load := map[int]float64{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		index, value, ok := strings.Cut(line, ",")
		if !ok {
			continue
		}
		i, err := strconv.Atoi(strings.TrimSpace(index))
		if err != nil {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if err != nil {
			continue
		}
		load[i] = v
	}
	return load
}

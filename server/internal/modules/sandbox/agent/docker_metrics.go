package agent

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type stats struct {
	CPU      float64
	NetBytes int64
}

func (d Docker) Stats(ctx context.Context, id string) (stats, error) {
	var raw struct {
		CPU    cpuStats `json:"cpu_stats"`
		PreCPU cpuStats `json:"precpu_stats"`
		Nets   map[string]struct {
			RX int64 `json:"rx_bytes"`
			TX int64 `json:"tx_bytes"`
		} `json:"networks"`
	}
	if err := d.call(ctx, http.MethodGet, "/containers/"+id+"/stats?stream=false", nil, &raw); err != nil {
		return stats{}, err
	}
	var out stats
	used := float64(raw.CPU.Usage.Total) - float64(raw.PreCPU.Usage.Total)
	system := float64(raw.CPU.System) - float64(raw.PreCPU.System)
	if used > 0 && system > 0 {
		out.CPU = used / system * float64(max(raw.CPU.Online, 1)) * 100
	}
	for _, n := range raw.Nets {
		out.NetBytes += n.RX + n.TX
	}
	return out, nil
}

type cpuStats struct {
	Usage struct {
		Total uint64 `json:"total_usage"`
	} `json:"cpu_usage"`
	System uint64 `json:"system_cpu_usage"`
	Online int    `json:"online_cpus"`
}

func (d Docker) Processes(ctx context.Context, id string) ([]string, error) {
	var top struct {
		Titles    []string   `json:"Titles"`
		Processes [][]string `json:"Processes"`
	}
	if err := d.call(ctx, http.MethodGet, "/containers/"+id+"/top?ps_args="+url.QueryEscape("-o pid,comm"), nil, &top); err != nil {
		return nil, err
	}
	column := len(top.Titles) - 1
	for i, t := range top.Titles {
		if t == "COMMAND" || t == "CMD" {
			column = i
		}
	}
	out := make([]string, 0, len(top.Processes))
	for _, p := range top.Processes {
		if column >= 0 && column < len(p) {
			if fields := strings.Fields(p[column]); len(fields) > 0 {
				out = append(out, fields[0])
			}
		}
	}
	return out, nil
}

func (d Docker) Logs(ctx context.Context, id string, lines int) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/containers/%s/logs?stdout=1&stderr=1&tail=%d", d.base, id, lines), nil)
	if err != nil {
		return "", err
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("docker logs: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("docker logs: %d", resp.StatusCode)
	}
	return demux(raw), nil
}

func demux(raw []byte) string {
	var b strings.Builder
	for len(raw) >= 8 && raw[0] <= 2 && raw[1] == 0 && raw[2] == 0 && raw[3] == 0 {
		size := int(binary.BigEndian.Uint32(raw[4:8]))
		raw = raw[8:]
		size = min(size, len(raw))
		b.Write(raw[:size])
		raw = raw[size:]
	}
	b.Write(raw)
	return b.String()
}

func (d Docker) Pause(ctx context.Context, id string) error {
	return d.call(ctx, http.MethodPost, "/containers/"+id+"/pause", nil, nil)
}

func (d Docker) Unpause(ctx context.Context, id string) error {
	return d.call(ctx, http.MethodPost, "/containers/"+id+"/unpause", nil, nil)
}

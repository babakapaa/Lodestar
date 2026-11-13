package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const apiVersion = "/v1.43"

type Docker struct {
	base   string
	client *http.Client
}

func NewDocker(host string) (Docker, error) {
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}
	if socket, ok := strings.CutPrefix(host, "unix://"); ok {
		transport := &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socket)
		}}
		return Docker{base: "http://docker" + apiVersion, client: &http.Client{Transport: transport}}, nil
	}
	if address, ok := strings.CutPrefix(host, "tcp://"); ok {
		return Docker{base: "http://" + address + apiVersion, client: http.DefaultClient}, nil
	}
	if strings.HasPrefix(host, "http://") {
		return Docker{base: strings.TrimSuffix(host, "/") + apiVersion, client: http.DefaultClient}, nil
	}
	return Docker{}, fmt.Errorf("unsupported DOCKER_HOST %q", host)
}

func (d Docker) List(ctx context.Context) ([]Container, error) {
	filters, err := json.Marshal(map[string][]string{"label": {labelRun}})
	if err != nil {
		return nil, err
	}
	var items []struct {
		ID     string            `json:"Id"`
		State  string            `json:"State"`
		Labels map[string]string `json:"Labels"`
	}
	if err := d.call(ctx, http.MethodGet, "/containers/json?all=true&filters="+url.QueryEscape(string(filters)), nil, &items); err != nil {
		return nil, err
	}
	out := make([]Container, 0, len(items))
	for _, item := range items {
		c := fromLabels(item.ID, item.State, item.Labels)
		if err := d.inspect(ctx, &c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (d Docker) inspect(ctx context.Context, c *Container) error {
	var info struct {
		State struct {
			ExitCode int `json:"ExitCode"`
		} `json:"State"`
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := d.call(ctx, http.MethodGet, "/containers/"+c.ID+"/json", nil, &info); err != nil {
		return err
	}
	c.ExitCode = info.State.ExitCode
	if p, ok := ports[c.Kind]; ok {
		for _, binding := range info.NetworkSettings.Ports[p] {
			if _, err := fmt.Sscan(binding.HostPort, &c.Port); err == nil {
				break
			}
		}
	}
	return nil
}

func (d Docker) Pull(ctx context.Context, image string) error {
	return d.call(ctx, http.MethodPost, "/images/create?fromImage="+url.QueryEscape(image), nil, nil)
}

func (d Docker) Create(ctx context.Context, name string, spec map[string]any) (string, error) {
	var created struct {
		ID string `json:"Id"`
	}
	err := d.call(ctx, http.MethodPost, "/containers/create?name="+url.QueryEscape(name), spec, &created)
	return created.ID, err
}

func (d Docker) Start(ctx context.Context, id string) error {
	return d.call(ctx, http.MethodPost, "/containers/"+id+"/start", nil, nil)
}

func (d Docker) Stop(ctx context.Context, id string) error {
	return d.call(ctx, http.MethodPost, "/containers/"+id+"/stop?t=30", nil, nil)
}

func (d Docker) Remove(ctx context.Context, id string) error {
	return d.call(ctx, http.MethodDelete, "/containers/"+id+"?force=true", nil, nil)
}

func (d Docker) call(ctx context.Context, method, path string, body, out any) error {
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, d.base+path, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("docker %s %s: %w", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	if resp.StatusCode >= 300 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("docker %s %s: %d %s", method, path, resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	if out == nil {
		_, err = io.Copy(io.Discard, resp.Body)
		return err
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

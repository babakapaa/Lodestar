package agent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/babakapaa/lodestar/server/internal/modules/sandbox"
)

var ports = map[string]string{
	sandbox.RunJupyter: "8888/tcp",
	sandbox.RunVSCode:  "3000/tcp",
	sandbox.RunDemo:    "8000/tcp",
}

func command(spec sandbox.RunSpec, token string) []string {
	switch spec.Kind {
	case sandbox.RunJupyter:
		return []string{"jupyter", "lab", "--ip=0.0.0.0", "--port=8888", "--no-browser", "--allow-root",
			"--ServerApp.token=" + token, "--ServerApp.root_dir=/workspace"}
	case sandbox.RunVSCode:
		return []string{"openvscode-server", "--host", "0.0.0.0", "--port", "3000", "--connection-token", token, "--default-folder", "/workspace"}
	default:
		return []string{"sh", "-c", spec.Command}
	}
}

func link(scheme, host string, c Container) string {
	if c.Port == 0 {
		return ""
	}
	base := fmt.Sprintf("%s://%s:%d", scheme, host, c.Port)
	switch c.Kind {
	case sandbox.RunJupyter:
		return base + "/lab?token=" + c.Token
	case sandbox.RunVSCode:
		return base + "/?tkn=" + c.Token
	case sandbox.RunDemo:
		return base + "/"
	default:
		return ""
	}
}

func containerSpec(spec sandbox.RunSpec, gpu int, token, workspace string) map[string]any {
	host := map[string]any{
		"Binds":       []string{workspace + ":/workspace"},
		"ShmSize":     2 << 30,
		"SecurityOpt": []string{"no-new-privileges"},
		"CapDrop":     []string{"NET_RAW", "SYS_ADMIN", "SYS_PTRACE", "MKNOD", "AUDIT_WRITE"},
		"PidsLimit":   4096,
		"DeviceRequests": []map[string]any{{
			"Driver": "nvidia", "DeviceIDs": []string{strconv.Itoa(gpu)}, "Capabilities": [][]string{{"gpu"}},
		}},
	}
	out := map[string]any{
		"Image":      spec.Image,
		"Cmd":        command(spec, token),
		"WorkingDir": "/workspace",
		"Env":        []string{"LODESTAR_RUN=" + spec.RunID},
		"Labels": map[string]string{
			labelRun: spec.RunID, labelGPU: strconv.Itoa(gpu), labelToken: token, labelKind: spec.Kind,
		},
		"HostConfig": host,
	}
	if p, ok := ports[spec.Kind]; ok {
		out["ExposedPorts"] = map[string]any{p: map[string]any{}}
		host["PortBindings"] = map[string]any{p: []map[string]string{{"HostPort": ""}}}
	}
	return out
}

func newToken() (string, error) {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}

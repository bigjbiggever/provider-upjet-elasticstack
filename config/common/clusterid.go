package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// ClusterScopedID returns "<cluster_uuid>/<externalName>" unless externalName
// already contains a '/'.
func ClusterScopedID(ctx context.Context, externalName string, terraformProviderConfig map[string]any) (string, error) {
	if strings.Contains(externalName, "/") {
		return externalName, nil
	}
	clusterUUID, err := GetClusterUUID(ctx, terraformProviderConfig)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", clusterUUID, externalName), nil
}

// ExternalNameFromStateID strips "<cluster_uuid>/" and returns only the
// resource identifier.
func ExternalNameFromStateID(tfstate map[string]any) (string, error) {
	id, ok := tfstate["id"].(string)
	if !ok || id == "" {
		return "", fmt.Errorf("cannot find id in tfstate")
	}
	if parts := strings.SplitN(id, "/", 2); len(parts) == 2 {
		return parts[1], nil
	}
	return id, nil
}

// GetClusterUUID discovers the cluster UUID using provider configuration.
//
//nolint:gocyclo // The branching is intentional to support multiple provider-config map shapes.
func GetClusterUUID(ctx context.Context, terraformProviderConfig map[string]any) (string, error) {
	cfg, err := resolveElasticsearchConfig(terraformProviderConfig)
	if err != nil {
		return "", err
	}

	endpoints, err := stringList(cfg["endpoints"])
	if err != nil || len(endpoints) == 0 {
		return "", fmt.Errorf("cannot resolve elasticsearch endpoints from provider config")
	}
	endpoint := strings.TrimRight(endpoints[0], "/")

	client := &http.Client{Timeout: 10 * time.Second}
	if insecure, ok := cfg["insecure"].(bool); ok && insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true}, //nolint:gosec
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/", nil)
	if err != nil {
		return "", err
	}
	if username, ok := cfg["username"].(string); ok && username != "" {
		if password, ok := cfg["password"].(string); ok && password != "" {
			req.SetBasicAuth(username, password)
		}
	}
	if token, ok := cfg["bearer_token"].(string); ok && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("failed to discover cluster UUID from %q: status=%d body=%q", endpoint, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	payload := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	clusterUUID, ok := payload["cluster_uuid"].(string)
	if !ok || clusterUUID == "" {
		return "", fmt.Errorf("cluster_uuid not found in elasticsearch response")
	}
	return clusterUUID, nil
}

func resolveElasticsearchConfig(terraformProviderConfig map[string]any) (map[string]any, error) {
	queue := []map[string]any{terraformProviderConfig}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if rawES, ok := current["elasticsearch"]; ok {
			return normalizeElasticsearchConfig(rawES)
		}
		for _, v := range current {
			if m, ok := asStringAnyMap(v); ok {
				queue = append(queue, m)
				continue
			}
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
				for i := 0; i < rv.Len(); i++ {
					if m, ok := asStringAnyMap(rv.Index(i).Interface()); ok {
						queue = append(queue, m)
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("terraform provider config does not include elasticsearch settings (root keys: %v)", mapKeys(terraformProviderConfig))
}

func normalizeElasticsearchConfig(rawES any) (map[string]any, error) {
	if cfg, ok := asStringAnyMap(rawES); ok {
		return cfg, nil
	}
	rv := reflect.ValueOf(rawES)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		if rv.Len() == 0 {
			return nil, fmt.Errorf("terraform provider config does not include elasticsearch settings")
		}
		first := rv.Index(0).Interface()
		cfg, ok := asStringAnyMap(first)
		if !ok {
			return nil, fmt.Errorf("invalid elasticsearch provider configuration item type %T", first)
		}
		return cfg, nil
	}
	return nil, fmt.Errorf("invalid elasticsearch provider configuration type %T", rawES)
}

func asStringAnyMap(v any) (map[string]any, bool) {
	if m, ok := v.(map[string]any); ok {
		return m, true
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map || rv.Type().Key().Kind() != reflect.String {
		return nil, false
	}
	out := make(map[string]any, rv.Len())
	iter := rv.MapRange()
	for iter.Next() {
		out[iter.Key().String()] = iter.Value().Interface()
	}
	return out, true
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func stringList(v any) ([]string, error) {
	switch s := v.(type) {
	case []string:
		return s, nil
	case []any:
		out := make([]string, 0, len(s))
		for _, e := range s {
			str, ok := e.(string)
			if !ok || str == "" {
				continue
			}
			out = append(out, str)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported endpoints type %T", v)
	}
}

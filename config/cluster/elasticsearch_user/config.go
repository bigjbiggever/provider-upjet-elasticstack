package elasticsearch_user

import (
	"context"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/bigjbiggever/provider-elasticstack/config/common"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("elasticstack_elasticsearch_security_user", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "elasticsearch"
		r.ShortGroup = "security"
		r.Kind = "ElasticsearchUser"
		// Keep username in spec.forProvider and derive full Terraform ID as
		// <cluster_uuid>/<username> at runtime so users do not need to provide UUID.
		r.ExternalName = config.NewExternalNameFrom(config.IdentifierFromProvider,
			config.WithGetIDFn(func(_ config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, terraformProviderConfig map[string]any) (string, error) {
				return getUserID(ctx, externalName, parameters, terraformProviderConfig)
			}),
			config.WithGetExternalNameFn(func(_ config.GetExternalNameFn, tfstate map[string]any) (string, error) {
				id, ok := tfstate["id"].(string)
				if !ok || id == "" {
					return "", fmt.Errorf("cannot find id in tfstate")
				}
				if parts := strings.SplitN(id, "/", 2); len(parts) == 2 {
					return parts[1], nil
				}
				return id, nil
			}),
		)
		r.ExternalName.OmittedFields = []string{"name"}
		delete(r.TerraformResource.Schema, "name")
		if s, ok := r.TerraformResource.Schema["username"]; ok {
			s.Optional = false
			s.Computed = false
		}
	})
}

func getUserID(ctx context.Context, externalName string, parameters map[string]any, terraformProviderConfig map[string]any) (string, error) {
	username := externalName
	if username == "" {
		if u, ok := parameters["username"].(string); ok {
			username = u
		}
	}
	if username == "" {
		return "", fmt.Errorf("cannot determine username from external name or parameters")
	}
	if strings.Contains(username, "/") {
		return username, nil
	}
	return common.ClusterScopedID(ctx, username, terraformProviderConfig)
}

package elasticsearch_role

import (
	"context"
	"fmt"
	"strings"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/bigjbiggever/provider-elasticstack/config/common"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("elasticstack_elasticsearch_security_role", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "github"
		r.ShortGroup = "security"
		r.Kind = "ElasticsearchRole"
		r.ExternalName = config.NewExternalNameFrom(config.NameAsIdentifier,
			config.WithGetIDFn(func(_ config.GetIDFn, ctx context.Context, externalName string, _ map[string]any, terraformProviderConfig map[string]any) (string, error) {
				return getRoleID(ctx, externalName, terraformProviderConfig)
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
	})
}

func getRoleID(ctx context.Context, externalName string, terraformProviderConfig map[string]any) (string, error) {
	return common.ClusterScopedID(ctx, externalName, terraformProviderConfig)
}

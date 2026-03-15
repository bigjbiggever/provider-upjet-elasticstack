package index_lifecycle

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/bigjbiggever/provider-elasticstack/config/common"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("elasticstack_elasticsearch_index_lifecycle", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "github"
		r.ShortGroup = "index"
		r.Kind = "IndexLifecycle"
		r.ExternalName = config.NewExternalNameFrom(config.NameAsIdentifier,
			config.WithGetIDFn(func(_ config.GetIDFn, ctx context.Context, externalName string, _ map[string]any, terraformProviderConfig map[string]any) (string, error) {
				return common.ClusterScopedID(ctx, externalName, terraformProviderConfig)
			}),
			config.WithGetExternalNameFn(func(_ config.GetExternalNameFn, tfstate map[string]any) (string, error) {
				return common.ExternalNameFromStateID(tfstate)
			}),
		)
	})
}

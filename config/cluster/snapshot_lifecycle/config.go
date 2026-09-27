package snapshot_lifecycle

import (
	"context"

	"github.com/crossplane/upjet/v2/pkg/config"

	"github.com/bigjbiggever/provider-elasticstack/config/common"
)

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("elasticstack_elasticsearch_snapshot_lifecycle", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "github"
		r.ShortGroup = "snapshot"
		r.Kind = "SnapshotLifecycle"
    r.ExternalName = config.NewExternalNameFrom(config.IdentifierFromProvider,
            config.WithGetIDFn(func(_ config.GetIDFn, ctx context.Context, externalName string, parameters map[string]any, terraformProviderConfig map[string]any) (string, error) {
                if externalName == "" {
                    if v, ok := parameters["name"]; ok {
                        if s, ok2 := v.(string); ok2 && s != "" {
                            externalName = s
                        }
                    }
                }
                return common.ClusterScopedID(ctx, externalName, terraformProviderConfig)
            }),
            config.WithGetExternalNameFn(func(_ config.GetExternalNameFn, tfstate map[string]any) (string, error) {
                return common.ExternalNameFromStateID(tfstate)
            }),
		)
	})
}

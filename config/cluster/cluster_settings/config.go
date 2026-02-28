package cluster_settings

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("elasticstack_elasticsearch_cluster_settings", func(r *config.Resource) {
		r.ShortGroup = "cluster"
		r.Kind = "ClusterSettings"
		r.ExternalName = config.IdentifierFromProvider
	})
}

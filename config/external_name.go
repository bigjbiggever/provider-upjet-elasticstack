package config

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// ExternalNameConfigs contains all external name configurations for this
// provider.
var ExternalNameConfigs = map[string]config.ExternalName{
    // Import requires using a randomly generated ID from provider: nl-2e21sda
    "elasticstack_elasticsearch_cluster_settings":   config.IdentifierFromProvider,
    // Allow K8s object name to differ from Elastic internal name.
    // Keep the "name" argument in the schema and let users set it explicitly.
    // Using IdentifierFromProvider here disables the name initializer and
    // prevents coupling metadata.name to the provider identifier.
    "elasticstack_elasticsearch_index_lifecycle":     config.IdentifierFromProvider,
    "elasticstack_elasticsearch_security_role":       config.IdentifierFromProvider,
    "elasticstack_elasticsearch_security_user":       config.IdentifierFromProvider,
    "elasticstack_elasticsearch_snapshot_lifecycle":  config.IdentifierFromProvider,
    "elasticstack_elasticsearch_snapshot_repository": config.IdentifierFromProvider,
}

// ExternalNameConfigurations applies all external name configs listed in the
// table ExternalNameConfigs and sets the version of those resources to v1beta1
// assuming they will be tested.
func ExternalNameConfigurations() config.ResourceOption {
	return func(r *config.Resource) {
		if e, ok := ExternalNameConfigs[r.Name]; ok {
			r.ExternalName = e
		}
	}
}

// ExternalNameConfigured returns the list of all resources whose external name
// is configured manually.
func ExternalNameConfigured() []string {
	l := make([]string, len(ExternalNameConfigs))
	i := 0
	for name := range ExternalNameConfigs {
		// $ is added to match the exact string since the format is regex.
		l[i] = name + "$"
		i++
	}
	return l
}

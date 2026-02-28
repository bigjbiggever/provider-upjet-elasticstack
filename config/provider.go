package config

import (
	// Note(turkenh): we are importing this to embed provider schema document
	_ "embed"

	ujconfig "github.com/crossplane/upjet/v2/pkg/config"

	clusterSettingsCluster "github.com/bigjbiggever/provider-elasticstack/config/cluster/cluster_settings"
	elasticsearchRoleCluster "github.com/bigjbiggever/provider-elasticstack/config/cluster/elasticsearch_role"
	elasticsearchUserCluster "github.com/bigjbiggever/provider-elasticstack/config/cluster/elasticsearch_user"
	indexLifecycleCluster "github.com/bigjbiggever/provider-elasticstack/config/cluster/index_lifecycle"
	snapshotLifecycleCluster "github.com/bigjbiggever/provider-elasticstack/config/cluster/snapshot_lifecycle"
	snapshotRepositoryCluster "github.com/bigjbiggever/provider-elasticstack/config/cluster/snapshot_repository"
	clusterSettingsNamespaced "github.com/bigjbiggever/provider-elasticstack/config/namespaced/cluster_settings"
	elasticsearchRoleNamespaced "github.com/bigjbiggever/provider-elasticstack/config/namespaced/elasticsearch_role"
	elasticsearchUserNamespaced "github.com/bigjbiggever/provider-elasticstack/config/namespaced/elasticsearch_user"
	indexLifecycleNamespaced "github.com/bigjbiggever/provider-elasticstack/config/namespaced/index_lifecycle"
	snapshotLifecycleNamespaced "github.com/bigjbiggever/provider-elasticstack/config/namespaced/snapshot_lifecycle"
	snapshotRepositoryNamespaced "github.com/bigjbiggever/provider-elasticstack/config/namespaced/snapshot_repository"
)

const (
	resourcePrefix = "elasticstack"
	modulePath     = "github.com/bigjbiggever/provider-elasticstack"
)

//go:embed schema.json
var providerSchema string

//go:embed provider-metadata.yaml
var providerMetadata string

// GetProvider returns provider configuration
func GetProvider() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("elasticstack.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		clusterSettingsCluster.Configure,
		indexLifecycleCluster.Configure,
		snapshotLifecycleCluster.Configure,
		snapshotRepositoryCluster.Configure,
		elasticsearchUserCluster.Configure,
		elasticsearchRoleCluster.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

// GetProviderNamespaced returns the namespaced provider configuration
func GetProviderNamespaced() *ujconfig.Provider {
	pc := ujconfig.NewProvider([]byte(providerSchema), resourcePrefix, modulePath, []byte(providerMetadata),
		ujconfig.WithRootGroup("elasticstack.m.crossplane.io"),
		ujconfig.WithIncludeList(ExternalNameConfigured()),
		ujconfig.WithFeaturesPackage("internal/features"),
		ujconfig.WithDefaultResourceOptions(
			ExternalNameConfigurations(),
		),
		ujconfig.WithExampleManifestConfiguration(ujconfig.ExampleManifestConfiguration{
			ManagedResourceNamespace: "crossplane-system",
		}))

	for _, configure := range []func(provider *ujconfig.Provider){
		// add custom config functions
		clusterSettingsNamespaced.Configure,
		indexLifecycleNamespaced.Configure,
		snapshotLifecycleNamespaced.Configure,
		snapshotRepositoryNamespaced.Configure,
		elasticsearchUserNamespaced.Configure,
		elasticsearchRoleNamespaced.Configure,
	} {
		configure(pc)
	}

	pc.ConfigureResources()
	return pc
}

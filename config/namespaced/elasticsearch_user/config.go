package elasticsearch_user

import "github.com/crossplane/upjet/v2/pkg/config"

// Configure configures individual resources by adding custom ResourceConfigurators.
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("elasticstack_elasticsearch_security_user", func(r *config.Resource) {
		// We need to override the default group that upjet generated for
		// this resource, which would be "elasticsearch"
		r.ShortGroup = "security"
		r.Kind = "ElasticsearchUser"
		// Keep username in spec.forProvider and let provider determine the external name.
		r.ExternalName = config.IdentifierFromProvider
		r.ExternalName.OmittedFields = []string{"name"}
		delete(r.TerraformResource.Schema, "name")
		if s, ok := r.TerraformResource.Schema["username"]; ok {
			s.Optional = false
			s.Computed = false
		}
	})
}

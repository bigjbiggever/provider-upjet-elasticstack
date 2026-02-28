package clients

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/v2/pkg/terraform"

	clusterv1beta1 "github.com/bigjbiggever/provider-elasticstack/apis/cluster/v1beta1"
	namespacedv1beta1 "github.com/bigjbiggever/provider-elasticstack/apis/namespaced/v1beta1"
)

const (
	keyUsername  = "username"
	keyPassword  = "password"
	keyEndpoints = "endpoints"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal elasticstack credentials as JSON"
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		pcSpec, err := resolveProviderConfig(ctx, client, mg)
		if err != nil {
			return terraform.Setup{}, errors.Wrap(err, "cannot resolve provider config")
		}

		data, err := resource.CommonCredentialExtractor(ctx, pcSpec.Credentials.Source, client, pcSpec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}

		// set provider configuration
		ps.Configuration = map[string]any{}
		es := map[string]any{}
		if v, ok := creds[keyUsername]; ok && v != "" {
			es[keyUsername] = v
		}
		if v, ok := creds[keyPassword]; ok && v != "" {
			es[keyPassword] = v
		}
		if v, ok := creds[keyEndpoints]; ok && v != "" {
			parts := strings.Split(v, ",")
			endpoints := make([]string, 0, len(parts))
			for _, p := range parts {
				s := strings.TrimSpace(p)
				if s != "" {
					endpoints = append(endpoints, s)
				}
			}
			if len(endpoints) > 0 {
				es[keyEndpoints] = endpoints
			}
		}
		if len(es) > 0 {
			ps.Configuration["elasticsearch"] = []map[string]any{es}
		}
		return ps, nil
	}
}

func resolveProviderConfig(ctx context.Context, crClient client.Client, mg resource.Managed) (*namespacedv1beta1.ProviderConfigSpec, error) {
	managed, ok := mg.(resource.ModernManaged)
	if !ok {
		return nil, errors.New("resource is not a managed resource")
	}
	return resolveModern(ctx, crClient, managed)
}

func resolveModern(ctx context.Context, crClient client.Client, mg resource.ModernManaged) (*namespacedv1beta1.ProviderConfigSpec, error) {
	configRef := mg.GetProviderConfigReference()
	if configRef == nil {
		return nil, errors.New(errNoProviderConfig)
	}

	candidates := []struct {
		gv               schema.GroupVersion
		pcu              resource.TypedProviderConfigUsage
		defaultNamespace string
	}{
		{
			gv:               clusterv1beta1.SchemeGroupVersion,
			pcu:              &clusterv1beta1.ProviderConfigUsage{},
			defaultNamespace: "",
		},
		{
			gv:               namespacedv1beta1.SchemeGroupVersion,
			pcu:              &namespacedv1beta1.ProviderConfigUsage{},
			defaultNamespace: mg.GetNamespace(),
		},
	}
	if mg.GetNamespace() != "" {
		candidates[0], candidates[1] = candidates[1], candidates[0]
	}

	for i := range candidates {
		c := candidates[i]
		pcRuntimeObj, err := crClient.Scheme().New(c.gv.WithKind(configRef.Kind))
		if err != nil {
			continue
		}
		pcObj, ok := pcRuntimeObj.(client.Object)
		if !ok {
			// This indicates a programming error, types are not properly generated.
			return nil, errors.New("provider config runtime object is not a client.Object")
		}

		lookupNamespace := c.defaultNamespace
		if configRef.Kind == namespacedv1beta1.ClusterProviderConfigKind || configRef.Kind == clusterv1beta1.ClusterProviderConfigKind {
			lookupNamespace = ""
		}

		if err := crClient.Get(ctx, types.NamespacedName{Name: configRef.Name, Namespace: lookupNamespace}, pcObj); err != nil {
			if apierrors.IsNotFound(err) {
				continue
			}
			return nil, errors.Wrap(err, errGetProviderConfig)
		}

		var pcSpec namespacedv1beta1.ProviderConfigSpec
		switch pc := pcObj.(type) {
		case *namespacedv1beta1.ProviderConfig:
			pcSpec = pc.Spec
			if pcSpec.Credentials.SecretRef != nil {
				pcSpec.Credentials.SecretRef.Namespace = mg.GetNamespace()
			}
		case *namespacedv1beta1.ClusterProviderConfig:
			pcSpec = pc.Spec
		case *clusterv1beta1.ProviderConfig:
			var convErr error
			pcSpec, convErr = convertClusterPCSpec(pc.Spec)
			if convErr != nil {
				return nil, errors.Wrap(convErr, "cannot convert cluster ProviderConfig spec")
			}
		case *clusterv1beta1.ClusterProviderConfig:
			var convErr error
			pcSpec, convErr = convertClusterPCSpec(pc.Spec)
			if convErr != nil {
				return nil, errors.Wrap(convErr, "cannot convert cluster ClusterProviderConfig spec")
			}
		default:
			return nil, errors.New("unknown provider config type")
		}

		t := resource.NewProviderConfigUsageTracker(crClient, c.pcu)
		if err := t.Track(ctx, mg); err != nil {
			return nil, errors.Wrap(err, errTrackUsage)
		}
		return &pcSpec, nil
	}
	return nil, errors.New(errGetProviderConfig)
}

func convertClusterPCSpec(pcSpec clusterv1beta1.ProviderConfigSpec) (namespacedv1beta1.ProviderConfigSpec, error) {
	data, err := json.Marshal(pcSpec)
	if err != nil {
		return namespacedv1beta1.ProviderConfigSpec{}, err
	}
	var shared namespacedv1beta1.ProviderConfigSpec
	if err := json.Unmarshal(data, &shared); err != nil {
		return namespacedv1beta1.ProviderConfigSpec{}, err
	}
	return shared, nil
}

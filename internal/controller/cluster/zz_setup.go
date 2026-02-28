// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	clustersettings "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/cluster/clustersettings"
	indexlifecycle "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/index/indexlifecycle"
	providerconfig "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/providerconfig"
	elasticsearchrole "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/security/elasticsearchrole"
	elasticsearchuser "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/security/elasticsearchuser"
	snapshotlifecycle "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/snapshot/snapshotlifecycle"
	snapshotrepository "github.com/bigjbiggever/provider-elasticstack/internal/controller/cluster/snapshot/snapshotrepository"
)

// Setup creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		clustersettings.Setup,
		indexlifecycle.Setup,
		providerconfig.Setup,
		elasticsearchrole.Setup,
		elasticsearchuser.Setup,
		snapshotlifecycle.Setup,
		snapshotrepository.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		clustersettings.SetupGated,
		indexlifecycle.SetupGated,
		providerconfig.SetupGated,
		elasticsearchrole.SetupGated,
		elasticsearchuser.SetupGated,
		snapshotlifecycle.SetupGated,
		snapshotrepository.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

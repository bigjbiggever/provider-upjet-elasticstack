# Provider Elastic Stack

`provider-elasticstack` is a [Crossplane](https://crossplane.io/) provider
elasticstack that is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the Elastic Stack
API.

## Getting Started

This elasticstack serves as a starting point for generating a new [Crossplane Provider](https://docs.crossplane.io/latest/packages/providers/) using the [`upjet`](https://github.com/crossplane/upjet) tooling. Please follow the guide linked below to generate a new Provider:

https://github.com/crossplane/upjet/blob/main/docs/generating-a-provider.md

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/bigjbiggever/provider-elasticstack/issues).

## NOTICE

This provider is not yet ready for production use. It is a work in progress and is not yet feature complete.

## Resources

| Resource | Group | Version | Kind | Status | Notes |
|----------|-------|---------|------|--------|-------|
| ElasticsearchClusterSettings | cluster | v1alpha1 | ClusterSettings | Working - Needs Testing | |
| ElasticsearchRole | security | v1alpha1 | ElasticsearchRole | Partially implemented - Not working | |
| ElasticsearchUser | security | v1alpha1 | ElasticsearchUser | Working - Needs Testing | |
| SnapshotLifecycle | snapshot | v1alpha1 | SnapshotLifecycle | Unknown | |
| SnapshotRepository | snapshot | v1alpha1 | SnapshotRepository | Unknown | |
| Index Lifecycle Management | index | v1alpha1 | IndexLifecycle | Unknown | |

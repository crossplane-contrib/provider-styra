# provider-styra

Crossplane provider for [Styra DAS](https://styra.com). The provider built from this repository can be installed into a Crossplane control plane or run seperately. It provides the following features:

* Extension of the K8s API with CRDs to represent Styra objects as K8s resources
* Controllers to provision these resources into a Styra instance
* Implementations of Crossplane's portable resource abstractions, enabling Styra resources to fulfill a user's general need for cloud services

## Getting Started and Documentation

For getting started guides, installation, deployment, and administration, see
our [Documentation](https://crossplane.io/docs/latest).

## Contributing

provider-styra is a community driven project and we welcome contributions. See the
Crossplane
[Contributing](https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md)
guidelines to get started.

### Adding New Resource

New resources can be added by defining the required types in `apis` and the controllers `pkg/controllers/`.

To generate the CRD YAML files run

    make generate


## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/crossplane-contrib/provider-styra/issues).

## Contact

Please use the following to reach members of the community:

* Slack: Join our [slack channel](https://slack.crossplane.io)
* Forums:
  [crossplane-dev](https://groups.google.com/forum/#!forum/crossplane-dev)
* Twitter: [@crossplane_io](https://twitter.com/crossplane_io)
* Email: [info@crossplane.io](mailto:info@crossplane.io)

## Governance and Owners

provider-aws is run according to the same
[Governance](https://github.com/crossplane/crossplane/blob/master/GOVERNANCE.md)
and [Ownership](https://github.com/crossplane/crossplane/blob/master/OWNERS.md)
structure as the core Crossplane project.

## Code of Conduct

provider-styra adheres to the same [Code of
Conduct](https://github.com/crossplane/crossplane/blob/master/CODE_OF_CONDUCT.md)
as the core Crossplane project.

## Licensing

provider-styra is under the Apache 2.0 license.


## Usage

To run the project

    make run

To run all tests:

    make test

To build the project

    make build

To list all available options

    make help

[See more](./INSTALL.md)

## Code generation

See [CODE_GENERATION.md](./CODE_GENERATION.md)

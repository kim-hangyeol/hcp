---
layout: ""
page_title: "Provider: Kubernetes alpha"
description: |-
  A generic provider for managing Kubernetes resources.
---

# Kubernetes alpha provider

This Kubernetes provider for Terraform (alpha) supports all API resources in a generic fashion.

This provider allows you to describe any Kubernetes resource using HCL. See [Moving from YAML to HCL](#moving-from-yaml-to-hcl) if you have YAML you want to use with the provider.

Please regard this project as experimental. It still requires extensive testing and polishing to mature into production-ready quality. Please [file issues](https://github.com/hashicorp/terraform-provider-kubernetes-alpha/issues/new/choose) generously and detail your experience while using the provider. We welcome your feedback.

Our eventual goal is for this generic resource to become a part of our [official Kubernetes provider](https://github.com/hashicorp/terraform-provider-kubernetes) once it is supported by the Terraform Plugin SDK. However, this work is subject to signficant changes as we iterate towards that level of quality.

## Schema

### Optional

- **client_certificate** (String, Optional) (env-var: `KUBE_CLIENT_CERT_DATA`) PEM-encoded client TLS certificate (including intermediates, if any).
- **client_key** (String, Optional) (env-var: `KUBE_CLIENT_KEY_DATA`) PEM-encoded private key for the above certificate.
- **cluster_ca_certificate** (String, Optional) (env-var: `KUBE_CLUSTER_CA_CERT_DATA`) PEM-encoded CA TLS certificate (including intermediates, if any).
- **config_context** (String, Optional) (env-var: `KUBE_CTX`) Context to select from the loaded `kubeconfig` file.
- **config_context_cluster** (String, Optional) (env-var: `KUBE_CTX_CLUSTER`) Cluster entry to associate to the current context (from kubeconfig).
- **config_context_user** (String, Optional) (env-var: `KUBE_CTX_USER`) User entry to associate to the current context (from kubeconfig).
- **config_path** (String, Optional) (env-var: `KUBE_CONFIG_PATH`) Path to a `kubeconfig` file.
- **exec** (Object, Optional) (see [below for nested schema](#nestedatt--exec))
- **host** (String, Optional) (env-var: `KUBE_HOST`) URL to the base of the API server.
- **insecure** (Boolean, Optional) (env-var: `KUBE_INSECURE`) Disregard invalid TLS certificates _(default false)_.
- **password** (String, Optional) (env-var: `KUBE_PASSWORD`) Basic authentication password.
- **token** (String, Optional) (env-var: `KUBE_TOKEN`) Token is a bearer token used by the client for request authentication.
- **username** (String, Optional) (env-var: `KUBE_USERNAME`) Basic authentication username.

<a id="nestedatt--exec"></a>
### Nested Schema for `exec`

- **api_version** (String) Version of the "client.authentication.k8s.io" API which the plugin implements.
- **args** (List of String) Command line arguments to the plugin command.
- **command** (String) The plugin executable (absolute path, or expects the plugin to be in OS PATH).
- **env** (Map of String) Environment values to set on the plugin process.

All attributes are optional, but you must either set a config path or static credentials. An empty provider block will not be a functional configuration.

Due to the internal design of this provider, access to a responsive API server is required both during PLAN and APPLY. The provider makes calls to the Kubernetes API to retrieve metadata and type information during all stages of Terraform operations.

### Credentials

For authentication, the provider can be configured with identity credentials sourced from either a `kubeconfig` file, explicit values in the `provider` block, or a combination of both.

If the `config_path` attribute is set to the path of a `kubeconfig` file, the provider will load it and use the credential values in it. When `config_path` is not set **NO EXTERNAL KUBECONFIG WILL BE LOADED**. Specifically, $KUBECONFIG environment variable is **NOT** considered.

Take note of the `current-context` configured in the file. You can override it using the `config_context` provider attribute.

If both `kubeconfig` and static credentials are defined in the `provider` block, the provider will prefer any attributes specified by the static credentials and ignore the corresponding attributes in the `kubeconfig`.

There are five options for providing identity information to the provider for authentication purposes:

* a kubeconfig
* a client certificate & key pair
* a static token
* a username & password pair
* an authentication plugin, such as `oidc` or `exec` (see examples folder).

## Experimental Status

By using the software in this repository (the "Software"), you acknowledge that: (1) the Software is still in development, may change, and has not been released as a commercial product by HashiCorp and is not currently supported in any way by HashiCorp; (2) the Software is provided on an "as-is" basis, and may include bugs, errors, or other issues;  (3) the Software is NOT INTENDED FOR PRODUCTION USE, use of the Software may result in unexpected results, loss of data, or other unexpected results, and HashiCorp disclaims any and all liability resulting from use of the Software; and (4) HashiCorp reserves all rights to make all decisions about the features, functionality and commercial release (or non-release) of the Software, at any time and without any obligation or liability whatsoever.

## Getting Started

If this is your first time here, you can get an overview of the provider by reading our [introductory blog post](https://www.hashicorp.com/blog/deploy-any-resource-with-the-new-kubernetes-provider-for-hashicorp-terraform/).

Next, review the "Schema" section above to understand which configuration options are available. You can find the following examples and more in [our examples folder](https://github.com/hashicorp/terraform-provider-kubernetes-alpha/blob/master/examples/). Don't forget to run `terraform init` in your Terraform configuration directory to allow Terraform to detect the provider plugin.

### Create a Kubernetes ConfigMap
```hcl
provider "kubernetes-alpha" {
  config_path = "~/.kube/config" // path to kubeconfig
}

resource "kubernetes_manifest" "test-configmap" {
  provider = kubernetes-alpha

  manifest = {
    "apiVersion" = "v1"
    "kind" = "ConfigMap"
    "metadata" = {
      "name" = "test-config"
      "namespace" = "default"
    }
    "data" = {
      "foo" = "bar"
    }
  }
}
```

### Create a Kubernetes Custom Resource Definition

```hcl
provider "kubernetes-alpha" {
  config_path = "~/.kube/config" // path to kubeconfig
}

resource "kubernetes_manifest" "test-crd" {
  provider = kubernetes-alpha

  manifest = {
    apiVersion = "apiextensions.k8s.io/v1"
    kind = "CustomResourceDefinition"
    metadata = {
      name = "testcrds.hashicorp.com"
    }
    spec = {
      group = "hashicorp.com"
      names = {
        kind = "TestCrd"
        plural = "testcrds"
      }
      scope = "Namespaced"
      versions = [{
        name = "v1"
        served = true
        storage = true
        schema = {
          openAPIV3Schema = {
            type = "object"
            properties = {
              data = {
                type = "string"
              }
              refs = {
                type = "number"
              }
            }
          }
        }
      }]
    }
  }
}
```

## Using `wait_for` to block create and update calls

The `kubernetes_manifest` resource supports the ability to block create and update calls until a field is set or has a particular value by specifying the `wait_for` attribute. This is useful for when you create resources like Jobs and Services when you want to wait for something to happen after the resource is created by the API server before Terraform should consider the resource created.

`wait_for` currently supports a `fields` attribute which allows you specify a map of fields paths to regular expressions. You can also specify `*` if you just want to wait for a field to have any value.

```hcl
resource "kubernetes_manifest" "test" {
  provider = kubernetes-alpha

  manifest = {
    // ...
  }

  wait_for = {
    fields = {
      # Check the phase of a pod
      "status.phase" = "Running"

      # Check a container's status
      "status.containerStatuses.0.ready" = "true",

      # Check an ingress has an IP
      "status.loadBalancer.ingress.0.ip" = "^(\\d+(\\.|$)){4}"

      # Check the replica count of a Deployment
      "status.readyReplicas" = "2"
    }
  }
}

```

## Moving from YAML to HCL

The `manifest` attribute of the `kubernetes_manifest` resource accepts any arbitrary Kubernetes API object, using Terraform's [map](https://www.terraform.io/docs/configuration/expressions.html#map) syntax. If you have YAML you want to use with this provider, we recommend that you convert it to a map as an initial step and then manage that resource in Terraform, rather than using `yamldecode()` inside the resource block. 

You can quickly convert a single YAML file to an HCL map using this one liner:

```
echo 'yamldecode(file("test.yaml"))' | terraform console
```

Alternatively, there is also an experimental command line tool [tfk8s](https://github.com/jrhouston/tfk8s) you could use to convert Kubernetes YAML manifests into complete Terraform configurations.

## Contributing

We welcome your contribution. Please understand that the experimental nature of this repository means that contributing code may be a bit of a moving target. If you have an idea for an enhancement or bug fix, and want to take on the work yourself, please first [create an issue](https://github.com/hashicorp/terraform-provider-kubernetes-alpha/issues/new/choose) so that we can discuss the implementation with you before you proceed with the work.

You can review our [contribution guide](https://github.com/hashicorp/terraform-provider-kubernetes-alpha/blob/master/_about/CONTRIBUTING.md) to begin. You can also check out our [frequently asked questions](https://github.com/hashicorp/terraform-provider-kubernetes-alpha/blob/master/_about/FAQ.md).

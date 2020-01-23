// In this example, we deploy [Grafana](https://grafana.com) using a `Deployment` to Kubernetes.

// First, import `ksonnet-lib`, the Jsonnet Kubernetes library which provides us
// helper methods for creating all kinds of Kubernetes objects.
(import "ksonnet-util/kausal.libsonnet") +

// We are going to use a `Deployment`, from the `apps/v1` group:
{
  deployment: $.apps.v1.deploy.new(
    // `metadata.name` and `replicas` and `spec.template.containers` need to be
    // specified.
    name="grafana",
    replicas=1,

    // To add the actual container, use the `container.new` from `core/v1`. You
    // need to set the name and the Docker container image to use.
    containers=[
      $.core.v1.container.new(
        name="grafana",
        image="grafana/grafana"
      )
    ],
  )
}

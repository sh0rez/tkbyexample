// In this example, we configure a [cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/).

(import 'ksonnet-util/kausal.libsonnet') +

{
  local cronJob = $.batch.v1beta1.cronJob,
  local container = $.core.v1.container,

  // We're going to define our container in isolation. \
  // `::` hides this from output.
  // We'll reference it below.
  container::
    container.new('backup-job', 'my-org/my-image') +
    container.withEnvMap({
      CLUSTER: 'my-database',
      RETENTION: '7d',
    }),

  // Now to define the cronjob.
  cronjob:

    cronJob.new() +

    cronJob.mixin.metadata
    .withName('my-backup-job') +

    cronJob.mixin.metadata.withLabels({
      app: 'backup',
    }) +

    // This is how you define your schedule.
    cronJob.mixin.spec
    .withSchedule('0 0 * * *') +

    // We can customize our history limits like this.
    cronJob.mixin.spec
    .withSuccessfulJobsHistoryLimit(1) +
    cronJob.mixin.spec
    .withFailedJobsHistoryLimit(3) +

    // Optionally set a concurrency policy - 'Allow', 'Forbid', or 'Replace'. \
    // 'Allow' is the default.
    cronJob.mixin.spec
    .withConcurrencyPolicy('Forbid') +

    cronJob.mixin.spec.jobTemplate.spec
    .withBackoffLimit(3) +

    cronJob.mixin.spec.jobTemplate.spec
    .template.spec.withRestartPolicy('OnFailure') +

    // This is where we reference the container we defined above!
    cronJob.mixin.spec.jobTemplate.spec
    .template.spec.withContainers($.container),
}

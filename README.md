<a href="https://opensource.newrelic.com/oss-category/#community-project"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Community_Project.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Project.png"><img alt="New Relic Open Source community project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Project.png"></picture></a>

![GitHub forks](https://img.shields.io/github/forks/newrelic/newrelic-bitmovin-analytics?style=social)
![GitHub stars](https://img.shields.io/github/stars/newrelic/newrelic-bitmovin-analytics?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/newrelic/newrelic-bitmovin-analytics?style=social)

![GitHub all releases](https://img.shields.io/github/downloads/newrelic/newrelic-bitmovin-analytics/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/newrelic/newrelic-bitmovin-analytics)
![GitHub last commit](https://img.shields.io/github/last-commit/newrelic/newrelic-bitmovin-analytics)
![GitHub Release Date](https://img.shields.io/github/release-date/newrelic/newrelic-bitmovin-analytics)

![GitHub issues](https://img.shields.io/github/issues/newrelic/newrelic-bitmovin-analytics)
![GitHub issues closed](https://img.shields.io/github/issues-closed/newrelic/newrelic-bitmovin-analytics)
![GitHub pull requests](https://img.shields.io/github/issues-pr/newrelic/newrelic-bitmovin-analytics)
![GitHub pull requests closed](https://img.shields.io/github/issues-pr-closed/newrelic/newrelic-bitmovin-analytics)

# New Relic Bitmovin Analytics Integration

This integration collects Bitmovin Analytics data using the Bitmovin Analytics
API and exports the data to New Relic Metrics.

## Table of Contents

* [Getting Started](#getting-started)
   * [On-host Deployment](#on-host)
   * [Docker Deployment](#docker)
   * [AWS Lambda Deployment](#aws-lambda)
* [Usage](#usage)
   * [Command Line Options](#command-line-options)
   * [Configuration](#configuration)
      * [General configuration](#general-configuration)
      * [Pipeline configuration](#pipeline-configuration)
      * [Log configuration](#log-configuration)
      * [Query configuration](#query-configuration)

## Getting Started

The New Relic Bitmovin Analytics integration can be deployed directly on a host,
deployed as a Docker container, or deployed as an AWS Lambda function.

### On-host

The New Relic Bitmovin Analytics integration provides binaries for the following
host environments.

* [Linux x86](link to release)
* [Linux amd64](link to release)
* [Windows x86](link to release)
* [Windows amd64](link to release)

#### Deploy the integration on host

To run the Bitmovin Analytics integration on a host, perform the following
steps.

1. Download the appropriate binary from the [latest release](https://github.com/newrelic/newrelic-bitmovin-analytics/releases).
1. Create a directory named `configs` in the directory containing the binary.
1. Create a file named `config.yml` in the `configs` directory and copy the
   contents of [`configs/standard-config.yml`](./configs/standard-config.yml)
   into it.
1. Edit the `config.yml` to [configure](#configuration) the integration
   appropriately for your environment.
1. Using the directory containing the binary as the current working directory,
   execute the binary using the appropriate [Command Line Options](#command-line-options).

### Docker

A Docker image for the Bitmovin Analytics integration is available at
[https://hub.docker.com/r/newrelic/newrelic-bitmovin-analytics](https://hub.docker.com/r/newrelic/newrelic-bitmovin-analytics).
This image can be used to run the integration inside a Docker container
[directly](#run-directly-from-dockerhub), as a base image for
[building a custom image](#extend-the-base-image), or using
[the provided `Dockerfile`](./build/package/Dockerfile) to
[build a custom image](#build-a-custom-image).

#### Run directly from [DockerHub](https://hub.docker.com/)

The Bitmovin Analytics integration [Docker image](https://hub.docker.com/r/newrelic/newrelic-bitmovin-analytics)
can be run directly from [DockerHub](https://hub.docker.com/). In this scenario,
the integration will use the [standard configuration](./configs/standard-config.yml).
The `NEW_RELIC_LICENSE_KEY` and the Bitmovin credentials _must_ be specified
using environment variables.

To use a custom configuration, the [`config.yml`](#configyml) must be mapped
into the running container. It can be mapped using the default filename or using
a custom filename. In the case of the latter, the `--config_path` [command line option](#command-line-options)
must be specified with the custom filename. The `NEW_RELIC_LICENSE_KEY` and the
Bitmovin credentials can be specifed in the [`config.yml`](#configyml) or using
environment variables.

In both case, additional environment variables can be passed to the container
using `docker run` with the [`-e`, `--env`, or `--env-file` options](https://docs.docker.com/reference/cli/docker/container/run/#env)
for [configuration parameters](#configuration) that can be specified via
environment variables. See below for examples.

**Example 1: Using the standard configuration**

In the following example, the integration is run inside a Docker container using
the standard configuration. The New Relic License Key and Bitmovin credentials
are specified using environment variables. No command line argument are passed
to the integration.

```bash
docker run -t --rm --name new-relic-bitmovin \
   -e NEW_RELIC_LICENSE_KEY=$NEW_RELIC_LICENSE_KEY \
   -e BITMOVINAPIKEY=$BITMOVINAPIKEY \
   -e BITMOVINLICENSEKEY=$BITMOVINLICENSEKEY \
   -e BITMOVINTENANTORG=$BITMOVINTENANTORG \
   newrelic/newrelic-bitmovin-analytics
```

**Example 2: Using a custom configuration with the default filename**

In the following example, the file `config.yml` in the current directory on the
host system is mapped to the default location in the container
(`configs/config.yml`). As in the previous example, the New Relic License Key
and Bitmovin credentials are specified using environment variables. No command
line argument are passed to the integration.

```bash
docker run -t --rm --name new-relic-bitmovin \
   -v "$PWD/config.yml":/usr/src/app/configs/config.yml \
   -e NEW_RELIC_LICENSE_KEY=$NEW_RELIC_LICENSE_KEY \
   -e BITMOVINAPIKEY=$BITMOVINAPIKEY \
   -e BITMOVINLICENSEKEY=$BITMOVINLICENSEKEY \
   -e BITMOVINTENANTORG=$BITMOVINTENANTORG \
   newrelic/newrelic-bitmovin-analytics
```

**Example 3: Using a custom configuration with a custom filename**

In the following example, the file `config.yml` in the current directory on the
host system is mapped to the location `configs/my_custom_config.yml` in the
container and the `-f` [command line option](#command-line-options) is used to
specify the custom filename. The full path is not needed as `/usr/src/app` is
the working directory when the integration runs in the container. As in the
previous example, the New Relic License Key and Bitmovin credentials are
specified using environment variables.

```bash
docker run -t --rm --name newrelic-bitmovin \
   -v "$PWD/config.yml":/usr/src/app/configs/my_custom_config.yml \
   -e NEW_RELIC_LICENSE_KEY=$NEW_RELIC_LICENSE_KEY \
   -e BITMOVINAPIKEY=$BITMOVINAPIKEY \
   -e BITMOVINLICENSEKEY=$BITMOVINLICENSEKEY \
   -e BITMOVINTENANTORG=$BITMOVINTENANTORG \
   newrelic/newrelic-bitmovin-analytics \
   --config_path configs/my_custom_config.yml
```

**Example 4: Using additional environment variables**

In the following example, additional environment variables are passed to the
container to configure the integration. In this case, the Bitmovin metric prefix
is specified using the [`BITMOVINMETRICPREFIX`](#bitmovinmetricprefix) environment
variable and the [`RUNASSERVICE`](#runasservice) environment variable is
specified to configure the integration to run as a standalone service.

```bash
docker run -t --rm --name newrelic-bitmovin \
   -e NEW_RELIC_LICENSE_KEY=$NEW_RELIC_LICENSE_KEY \
   -e BITMOVINAPIKEY=$BITMOVINAPIKEY \
   -e BITMOVINLICENSEKEY=$BITMOVINLICENSEKEY \
   -e BITMOVINTENANTORG=$BITMOVINTENANTORG \
   -e BITMOVINMETRICPREFIX=my.bitmovin. \
   -e RUNASSERVICE=true \
   newrelic/newrelic-bitmovin-analytics
```

#### Extend the base image

The Bitmovin Analytics integration [Docker image](https://hub.docker.com/r/newrelic/newrelic-bitmovin-analytics)
can be used as the base image for building custom images. This scenario can be
easier as the [`config.yml`](#configyml) can be packaged into the custom image
and does not need to be mounted in. However, it does require access to
[a Docker registry](https://docs.docker.com/guides/docker-concepts/the-basics/what-is-a-registry/)
where the custom image can be pushed (e.g. [ECR](https://aws.amazon.com/ecr/)
and that is accessible to the technology used to manage the container
(e.g. [ECS](https://aws.amazon.com/ecs/)). In addition, this scenario requires
maintenance of a custom `Dockerfile` and the processes to build and publish the
image to a registry.

The minimal example of a `Dockerfile` for building a custom image simply extends
the base image (`newrelic/newrelic-bitmovin-analytics`) and copies a
configuration file to the default location (`/usr/src/app/config.yml`).

```dockerfile
FROM newrelic/newrelic-bitmovin-analytics

#
# Copy your config file into the default location.
# Adjust the local path as necessary.
#
COPY ./config.yml ./configs/config.yml
```

Note that the full directory path in the container does not need to be
specified. This is because the base image sets the
[`WORKDIR`](https://docs.docker.com/reference/dockerfile/#workdir) to
`/usr/src/app`. In fact, custom `Dockerfile`s should _not_ change the `WORKDIR`.

The following commands can be used to build a custom image using a custom
`Dockerfile` that extends the base image.

```bash
docker build -t newrelic-bitmovin-analytics-custom -f Dockerfile-custom .
docker tag newrelic-bitmovin-analytics-custom someregistry/username/newrelic-bitmovin-analytics-custom
docker push someregistry/username/newrelic-bitmovin-analytics-custom
```

Subsequently, the integration can be run using the custom image as in the
previous examples but without the need to mount the configuration file.

#### Build a custom image

The Bitmovin Analytics integration [Docker image](https://hub.docker.com/r/newrelic/newrelic-bitmovin-analytics)
can also be built locally using the provided [`Dockerfile`](./build/package/Dockerfile)
"as-is" or as the basis for building a custom `Dockerfile`. As is the case when
extending the base image, this scenario does require access to
[a Docker registry](https://docs.docker.com/guides/docker-concepts/the-basics/what-is-a-registry/)
where the custom image can be pushed (e.g. [ECR](https://aws.amazon.com/ecr/)
and that is accessible to the technology used to manage the container
(e.g. [ECS](https://aws.amazon.com/ecs/)). Similarly, it requires maintenance of
a custom `Dockerfile` and the processes to build and publish the image to a
registry.

The general procedure for building a custom image using the provided
`Dockerfile` "as-is" are as follows.

1. Clone this repository
1. Navigate to the repository root
1. Run the following commands

```bash
make docker
docker tag newrelic-bitmovin-analytics someregistry/username/newrelic-bitmovin-analytics
docker push someregistry/username/newrelic-bitmovin-analytics
```

To use a custom `Dockerfile`, backup the provided `Dockerfile`, make necessary
changes to the original, and follow the steps above.

As is the case when extending the base image, the integration can be run using
the custom image as in the previous examples but without the need to mount any
files into the container.

### AWS Lambda

The New Relic Bitmovin Analytics integration can be deployed and run as an AWS
Lambda function using [the provided CloudFormation template](./deployments/lambda/cf-template.yaml).
A [sample parameters file](./deployments/lambda/cf-params.sample.json) is
provided that shows an example of each parameter that can be used with the
CloudFormation template.

When used [independently](#deploy-using-the-makefile), the provided CloudFormation template will create a
new [AWS Stack](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacks.html)
with four resources.

1. The Lambda function
1. An [EventBridge](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-what-is.html)
   [schedule](https://docs.aws.amazon.com/eventbridge/latest/userguide/using-eventbridge-scheduler.html)
   that will invoke the Lambda function once a minute
1. An [EventBridge](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-what-is.html)
   schedule group to contain the schedule
1. An [IAM role](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html)
   that the [EventBridge](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-what-is.html)
   scheduler service can assume to execute the Lambda function

The only _required_ resource in this [Stack](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacks.html)
is the Lambda function. The [EventBridge](https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-what-is.html)
resources and the [IAM role](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html)
can be removed from the template if existing resources will be used and/or if
the Lambda function will be invoked as the target of another resource.

The same holds true when the Lambda function is deployed using
[customer-specific](#deploy-using-customer-specific-procedures) deployment
procedures. In addition, in this case, the Lambda function can be still be
deployed as it's own stack or it can be integrated as part of another stack.

#### Requirements for running the AWS Lambda

There are two requirements in order to run the integration as an AWS Lambda.

1. The AWS Lambda requires an execution role that the Lambda service can
   assume to run the Lambda function. Either an existing role can be used or
   a new role can be created.

   When creating a new role, the Lambda does not need access to any AWS
   services other than [CloudWatch](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/WhatIsCloudWatch.html)
   (to write log events). Attaching the
   [AWSLambdaExecute](https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AWSLambdaExecute.html)
   managed policy or the [AWSLambdaBasicExecutionRole](https://docs.aws.amazon.com/aws-managed-policy/latest/reference/AWSLambdaBasicExecutionRole.html)
   managed policy or the equivalent thereof, is sufficient.
1. The AWS Lambda requires an S3 bucket where the lambda deployment package
   can be uploaded to so that the [CloudFormation template](./deployments/lambda/cf-template.yaml)
   can reference it to during deployment. Either an existing bucket can be
   used or a new bucket  can be created.

#### Deploy using customer-specific procedures

Some organizations may have existing procedures for provisioning AWS resources
via CloudFormation, for example, as a part of a CI pipeline via tools like
[Terraform](https://www.terraform.io/) or [Ansible](https://docs.ansible.com/).
Use the [CloudFormation template](./deployments/lambda/cf-template.yaml)
and the [CloudFormation parameters](./deployments/lambda/cf-params.sample.json)
as a guide to deploy the integration as part of those procedures.

When using existing procedures for deployment, the Lambda
[deployment package](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)
can be prepared using the following steps.

1. Clone this repository
1. Navigate to the repository root
1. Run the command `make package-lambda`

On successful completion of the `make` process, the deployment package will be
located at `dist/newrelic-bitmovin-lambda.zip`. The deployment package will need
to be uploaded to the S3 bucket referenced in the [CloudFormation template](./deployments/lambda/cf-template.yaml).

#### Deploy using the [`Makefile`](./Makefile)

Alternately, the lambda can be deployed independently using the provided [`Makefile`](./Makefile).
This method requires the [AWS CLI](https://aws.amazon.com/cli/) to be installed
on the same system where the repository was cloned. To deploy the lambda using
this method, perform the following steps.

1. Clone this repository
1. Navigate to the repository root
1. Copy the [`configs/standard-config.yml`](./configs/standard-config.yml) to
   `configs/config.yml`
1. Copy the [`deployments/lambda/cf-params.sample.json`](./deployments/lambda/cf-params.sample.json)
   to `deployments/lambda/cf-params.json`
1. Use the parameter descriptions in the CloudFormation template at
   [`deployments/lambda/cf-template.yaml`](./deployments/lambda/cf-template.yaml)
   as a guide to update the template parameters in `deployments/lambda/cf-params.json`. \

   The only _required_ parameters are `ExecRoleArn`, `S3BucketName`, and the
   various New Relic and Bitmovin API/license keys. The others have sufficient
   defaults.
1. Ensure that appropriate [authentication and access credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-authentication.html)
   are set to allow the AWS CLI to authenticate.
1. Run the command `make deploy-lambda`
1. Verify that the command ran successfully by looking for the following console
   output.

```bash
...
Building lambda zip package...
updating: bootstrap (deflated 51%)
updating: configs/ (stored 0%)
updating: configs/standard-config.yml (deflated 54%)
updating: configs/config.yml (deflated 53%)
Uploading lambda zip package...
upload: ../newrelic-bitmovin-lambda.zip to s3://xxxxx/newrelic-bitmovin-lambda.zip
Deploying stack newrelic-bitmovin-analytics...

Waiting for changeset to be created..
Waiting for stack create/update to complete
Successfully created/updated stack - newrelic-bitmovin-analytics
Done.
```

## Usage

### Command Line Options

| Option | Description | Default |
| --- | --- | --- |
| --config_path | path to the (#configyml) to use | `configs/config.yml` |
| --dry_run | flag to enable "dry run" mode | `false` |
| --env_prefix | prefix to use for environment variable lookup | `''` |
| --verbose | flag to enable "verbose" mode | `false` |
| --version | display version information only | N/a |

### Configuration

#### `config.yml`

The Bitmovin Analytics integration is configured using a YAML file named
[`config.yml`](#configyml). The default location for this file is
`configs/config.yml` relative to the current working directory when the
integration binary is executed. The supported configuration
parameters are listed below. See [`standard-config.yml`](./config/standard-config.yml)
for a full configuration example.

##### General configuration

The parameters in this section are configured at the top level of the
[`configy.yml`](#configyml).

###### `licenseKey`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| New Relic license key | string | Y | N/a |

This parameter specifies the New Relic License Key (INGEST) that should be used
to send generated metrics.

The license key can also be specified using the `NEW_RELIC_LICENSE_KEY`
environment variable.

###### `region`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| New Relic region identifier | `US` / `EU` | N | `US` |

This parameter specifies which New Relic region that generated metrics should be
sent to.

###### `interval`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Polling interval (in _seconds_) | numeric | N | 60 |

This parameter specifies the interval (in _seconds_) at which the integration
should poll the Bitmovin Analytics API for metrics.

This parameter is only used when [`runAsService`](#runasservice) is set to
`true`.

###### `runAsService`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Flag to enable running the integration as a "service" | `true` / `false` | N | `false` |

The integration can run either as a "service" or as a simple command line
utility which runs once and exits when it is complete.

When set to `true`, the integration process will run continuously and poll the
Bitmovin Analytics API at the recurring interval specified by the [`interval`](#interval)
parameter. The process will only exit if it is explicitly stopped or a fatal
error or panic occurs.

When set to `false`, the integration will run once and exit. This is intended for
use with an external scheduling mechanism like [cron](https://man7.org/linux/man-pages/man8/cron.8.html).

###### `pipeline`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root element for the set of [pipeline configuration](#pipeline-configuration) parameters  | YAML Sequence | N | N/a |

The integration retrieves, processes, and exports metrics to New Relic using
a data pipeline consisting of one or more receivers, a processing chain, and a
New Relic exporter. Various aspects of the pipeline are configurable. This
element groups together the configuration parameters related to
[pipeline configuration](#pipeline-configuration).

###### `log`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The root node for the set of [log configuration](#log-configuration) parameters  | YAML Sequence | N | N/a |

The integration uses the [logrus](https://pkg.go.dev/github.com/sirupsen/logrus)
package for application logging. This element groups together the configuration
parameters related to [log configuration](#log-configuration).

###### `bitmovinMetricPrefix`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| A prefix to prepend to Bitmovin metric names | string | N | `''` |

This parameter specifies a prefix that will be preprended to each Bitmovin
metric name when the metric is exported to New Relic. Note that this is
_in addition_ to the automatically prepended prefixes added to distinguish
between the aggregation functions used for a query (e.g. `cnt`, `avg`, etc).

For example, if this parameter is set to `bitmovin.`, for a `count` query run
for the `IMPRESSION_ID` metric, the resulting metric name in New Relic would be
`bitmovin.cnt_impression_id`.

###### `bitmovinApiKey`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Bitmovin API key | string | Y | N/a |

This parameter specifies the [Bitmovin API key](https://developer.bitmovin.com/encoding/docs/get-started-with-the-bitmovin-api#get-your-bitmovin-api-key)
that will be used when making requests to the Bitmovin Analytics API.

The API key can also be specified using the `BITMOVINAPIKEY` environment
variable.

###### `bitmovinLicenseKey`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Bitmovin Analytics license key | string | Y | N/a |

This parameter specifies the [Bitmovin Analytics license key](https://developer.bitmovin.com/playback/docs/setup-analytics#analytics-licenses)
that will be used whe making requests to the Bitmovin Analytics API.

The license key can also be specified using the `BITMOVINLICENSEKEY` environment
variable.

###### `bitmovinTenantOrg`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Bitmovin tenant org ID | string | Y | N/a |

This parameter specifies the [Bitmovin tenant organization ID](https://developer.bitmovin.com/streams/docs/manage-organization-details)
that will be used whe making requests to the Bitmovin Analytics API.

The tenant org can also be specified using the `BITMOVINTENANTORG` environment
variable.

##### `queries`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| An array of [query](#query-configuration) configurations | YAML sequence | N | `[]` |

The Bitmovin Analytics metrics that the integration collects are configured
using the `queries` array. Each element of the array is a [query configuration](#query-configuration)
that defines the metric to collect, the filters to apply, the dimensions to
group by, and the criteria to order by.

See the [query configuration](#query-configuration) section for more details.

##### Pipeline configuration

###### `bufferSize`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Size of the buffer that holds items before processing | number | N | 500 |

This parameter specifies the size of the buffer that holds received items before
being flushed through the processing chain and on to the exporters. When this
size is reached, the items in the buffer will be flushed automatically.

###### `harvestTime`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Harvest interval (in _seconds_) | number | N | 60 |

This parameter specifies the interval (in _seconds_) at which the pipeline
should automatically flush received items through the processing chain and on
to the exporters. Each time this interval is reached, the pipeline will flush
items even if the item buffer has not reached the size specified by the
[`bufferSize`](#bufferSize) parameter.

###### `instances`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Number of concurrent pipeline instances to run | number | N | 3 |

The integration retrieves, processes, and exports metrics to New Relic using
a data pipeline consisting of one or more receivers, a processing chain, and a
New Relic exporter. When [`runAsService`](#runasservice) is `true`, the
integration can launch one or more "instances" of this pipeline to receive,
process, and export data concurrently. Each "instance" will be configured with
the same processing chain and exporter and the receivers will be spread across
the available instance in a round-robin fashion.

This parameter specifies the number of pipeline instances to launch.

**NOTE:** When [`runAsService`](#runasservice) is `false`, only a single
instance is used.

##### Log configuration

###### `level`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Log level | `panic` / `fatal` / `error` / `warn` / `info` / `debug` / `trace`  | N | `warn` |

This parameter specifies the maximum severity of log messages to output with
`trace` being the least severe and `panic` being the most severe. For example,
at the default log level (`warn`), all log messages with severities `warn`,
`error`, `fatal`, and `panic` will be output but `info`, `debug`, and `trace`
will not.

###### `fileName`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| Path to a file where log output will be written | string | N | `stderr` |

This parameter designates a file path where log output should be written. When
no path is specified, log output will be written to standard error (`stderr`).

##### Query configuration

The Bitmovin Analytics integration can collect metrics by querying the following
Bitmovin Analytics API endpoints.

* [Max Concurrent Viewers](https://developer.bitmovin.com/playback/reference/getanalyticsmetricsmaxconcurrentviewers)
* [Average Concurrent Viewers](https://developer.bitmovin.com/playback/reference/getanalyticsmetricsavgconcurrentviewers)
* [Average Dropped Frames](https://developer.bitmovin.com/playback/reference/getanalyticsmetricsavgconcurrentviewers)
* [Count Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriescount)
* [Sum Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriessum)
* [Average Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriesavg)
* [Min Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriesmin)
* [Max Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriesmax)
* [Standard Deviation Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriesstddev)
* [Percentile Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriespercentile)
* [Variance Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriesvariance)
* [Median Queries](https://developer.bitmovin.com/playback/reference/postanalyticsqueriesmedian)

The endpoints to call and the parameters used to call them are specified as
elements of the [`queries`](#queries) configuration parameter. Each element
_must_ contain the [`type`](#query-type) parameter and _may_ contain the
following additional parameters.

* [`metric`](#query-metric)
* [`dimensions`](#query-dimensions)
* [`filters`](#query-filters)
* [`orderBy`](#query-orderby)
* [`percentile`](#query-percentile)
* [`interval`](#query-interval)

###### Query `metric`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The Bitmovin Analytics API field name to query | Any Bitmovin Analytics [API field](https://developer.bitmovin.com/playback/docs/analytics-api-fields) | conditional | N/a |

This query configuration parameter specifies the Bitmovin Analytics
[API field](https://developer.bitmovin.com/playback/docs/analytics-api-fields)
to query, specified using upper snake case. For instance, the
API field `impressionId` would be specified as `IMPRESSION_ID`.

This configuration parameter is _required_ unless the [query `type`](#query-type)
is set to `max_concurrentviewers`, `avg_concurrentviewers`, or `avg_dropped_frames`.

###### Query `type`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The type of query to run | `max_concurrentviewers` / `avg_concurrentviewers` / `avg_dropped_frames` / `count` / `sum` / `average` / `min` / `max` / `stddev` / `percentile` / `variance` / `median` | Y | N/a |

This query configuration parameter specifies the Bitmovin Analytics API endpoint
to query. This parameter may have one of the following values.

| Value | API Reference |
| --- | --- |
| `max_concurrentviewers` | https://developer.bitmovin.com/playback/reference/getanalyticsmetricsmaxconcurrentviewers |
| `avg_concurrentviewers` | https://developer.bitmovin.com/playback/reference/getanalyticsmetricsavgconcurrentviewers |
| `avg_dropped_frames` | https://developer.bitmovin.com/playback/reference/getanalyticsmetricsavgdroppedframes |
| `count` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriescount |
| `sum` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriessum |
| `average` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriesavg |
| `min` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriesmin |
| `max` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriesmax |
| `stddev` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriesstddev |
| `percentile` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriespercentile |
| `variance` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriesvariance |
| `median` | https://developer.bitmovin.com/playback/reference/postanalyticsqueriesmedian |

See the [Query Examples section](#query-examples) for example usages.

###### Query `dimensions`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| An array of strings to group results by | YAML Sequence | N | `[]` |

This query configuration parameter is a list of strings, where each string is
the name of a Bitmovin Analytics [API field](https://developer.bitmovin.com/playback/docs/analytics-api-fields).
API field names are specified using upper snake case. For instance, the
following example would group results by `impressionId` and `cdnProvider`.

```yaml
dimensions:
- IMPRESSION_ID
- CDN_PROVIDER
```

See the [Query Examples section](#query-examples) for additional examples.

###### Query `filters`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| A set of criteria to filter results by | YAML Mapping | N | `{}` |

This parameter is a set of key/value pairs where the key is an
[API field](https://developer.bitmovin.com/playback/docs/analytics-api-fields)
name and the value is a set with exactly two key/value pairs: the filter
`operation` type and a constraint `value`. API field names are specified using
upper snake case. The `operation` key/value pair may have one of the following
values.

| Filter Operation |
| --- |
| `IN` |
| `EQ` |
| `NE` |
| `LT` |
| `LTE` |
| `GT` |
| `GTE` |
| `CONTAINS` |
| `NOTCONTAINS` |

For instance, the following example would filter results to
those where the `DURATION` value is greater than `100`.

```yaml
filters:
  DURATION:
    operation: GT
    value: 100
```

See the [Query Examples section](#query-examples) for additional examples.

###### Query `orderBy`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| An array of criteria to order results by | YAML Sequence | N | `[]` |

This parameter is a list of key/value pairs that specify an
[API field](https://developer.bitmovin.com/playback/docs/analytics-api-fields)
`name` and a sort `order`. API field names are specified using upper snake case.
The `order` may be the value `DESC` or `ASC`. For instance, the following
example would sort the query results first by `cdnProvider` in `DESC`ending
order and second by `deviceClass` in `ASC`ending order.

```yaml
orderBy:
- name: CDN_PROVIDER
  order: DESC
- name: DEVICE_CLASS
  order: ASC
```

See the [Query Examples section](#query-examples) for additional examples.

###### Query `percentile`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The percentile value to use for `percentile` queries | `0` - `99` | Y for `percentile` queries | N/a |

This parameter is a required value for `percentile` [query types](#query-type)
and specifies the percentile to use as a number from `0` to `99`.

###### Query `interval`

| Description | Valid Values | Required | Default |
| --- | --- | --- | --- |
| The timeseries granularity | `MINUTE` / `HOUR` / `DAY` / `MONTH` | N | `MINUTE` |

This parameter specifies the granularity of the query results. The following
values are supported.

| Granularity |
| --- |
| `MINUTE` |
| `HOUR` |
| `DAY` |
| `MONTH` |

See the [Query Examples section](#query-examples) for example usages.

###### Bitmovin to New Relic metric mapping

Metrics returned from the Bitmovin API are all mapped to
[New Relic gauge metrics](https://docs.newrelic.com/docs/data-apis/understand-data/metric-data/metric-data-type/)
and named using the [metric prefix](#bitmovinmetricprefix), the [query `type`](#query-type),
and the specified [query `metric`](#query-metric), converted to lower snake case. The
following table shows examples of the resulting New Relic metric names for each
[query `type`](#query-type), assuming the [metric prefix](#bitmovinmetricprefix)
is `bitmovin.`

| Query `type` | Query `metric` | New Relic Metric Name |
| --- | --- | --- |
| `max_concurrentviewers` | N/a | `bitmovin.max_concurrent_viewers` |
| `avg_concurrentviewers` | N/a | `bitmovin.avg_concurrent_viewers` |
| `avg_dropped_frames` | N/a | `bitmovin.avg_dropped_frames` |
| `count` | `IMPRESSION_ID` | `bitmovin.cnt_impression_id` |
| `sum` | `IMPRESSION_ID` | `bitmovin.sum_impression_id` |
| `average` | `IMPRESSION_ID` | `bitmovin.avg_impression_id` |
| `min` | `IMPRESSION_ID` | `bitmovin.min_impression_id` |
| `max` | `IMPRESSION_ID` | `bitmovin.max_impression_id` |
| `stddev` | `IMPRESSION_ID` | `bitmovin.stddev_impression_id` |
| `percentile` | `IMPRESSION_ID` | `bitmovin.pNN_impression_id` (where `NN` is the [query `percentile`](#query-percentile) ) |
| `variance` | `IMPRESSION_ID` | `bitmovin.var_impression_id` |
| `median` | `IMPRESSION_ID` | `bitmovin.med_impression_id` |

###### Query examples

The following examples show how to recreate each of the [Audience metrics](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#audience-metrics)
on the Bitmovin Analytics Dashboard.

**[Plays](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#audience-metrics)**

```yaml
queries:
# ...
- type: count
  metric: IMPRESSION_ID
  filters:
    VIDEO_STARTUPTIME:
      operation: GT
      value: 0
```

**[Play Attempts](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#play-attempts)**

```yaml
queries:
# ...
- type: count
  metric: PLAY_ATTEMPTS
```

**[Unique Users](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#unique-users)**

```yaml
queries:
# ...
- type: count
  metric: USER_ID
  filters:
    VIDEO_STARTUPTIME:
      operation: GT
      value: 0
```

**[Concurrent Viewers](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#concurrent-viewers)**

```yaml
queries:
# ...
- type: count
  metric: max_concurrentviewers
```

**[Total Page Loads](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#total-page-loads)**

```yaml
queries:
# ...
- type: count
  metric: IMPRESSION_ID
  filters:
    PLAYER_STARTUPTIME:
      operation: GT
      value: 0
```

**[Total Hours Watched](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#total-hours-watched)**

```yaml
queries:
# ...
- type: sum
  metric: PLAYED
  filters:
    PLAYED:
      operation: GT
      value: 0
```

**[Average View Time](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#average-view-time)**

```yaml
queries:
# ...
- type: average
  metric: VIEWTIME
```

The following examples show how to recreate each of the [Quality of Experience metrics](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#quality-of-experience)
on the Bitmovin Analytics Dashboard.

**[Total Startup Time](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#total-startup-time)**

```yaml
queries:
# ...
- type: median
  metric: STARTUPTIME
  filters:
    PAGE_LOAD_TYPE:
      operation: EQ
      value: 1
    STARTUPTIME:
      operation: GT
      value: 0
```

**[Player Startup Time](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#player-startup-time)**

```yaml
queries:
# ...
- type: median
  metric: PLAYER_STARTUPTIME
  filters:
    PAGE_LOAD_TYPE:
      operation: EQ
      value: 1
    STARTUPTIME:
      operation: GT
      value: 0
```

**[Video Startup Time](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#video-startup-time)**

```yaml
queries:
# ...
- type: median
  metric: VIDEO_STARTUPTIME
  filters:
    VIDEO_STARTUPTIME:
      operation: GT
      value: 0
```

**[Seek Time](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#seek-time)**

```yaml
queries:
# ...
- type: median
  metric: SEEKED
  filters:
    SEEKED:
      operation: GT
      value: 0
```

**[Error Percentage](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#error-percentage)**

```yaml
queries:
# ...
- type: average
  metric: ERROR_PERCENTAGE
```

**[Start Failures](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#start-failures)**

```yaml
queries:
# ...
- type: count
  metric: VIDEOSTART_FAILED
  filters:
    VIDEOSTART_FAILED_REASON:
      operation: NE
      value: PAGE_CLOSED
    VIDEOSTART_FAILED:
      operation: EQ
      value: true
```

**[Rebuffer Percentage](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#rebuffer-percentage)**

```yaml
queries:
# ...
- type: average
  metric: REBUFFER_PERCENTAGE
```

**[Buffering Time](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#buffering-time)**

```yaml
queries:
# ...
- type: average
  metric: BUFFERED
```

**[Data Downloaded](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#data-downloaded)**

```yaml
queries:
# ...
- type: sum
  metric: VIDEO_SEGMENTS_DOWNLOAD_SIZE
```

**[Bandwidth](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#bandwidth)**

```yaml
queries:
# ...
- type: average
  metric: DOWNLOAD_SPEED
```

**[Video Bitrate](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#video-bitrate)**

```yaml
queries:
# ...
- type: average
  metric: VIDEO_BITRATE
  filters:
    VIDEO_BITRATE:
      operation: GT
      value: 0
```

**[Scale Factor](https://developer.bitmovin.com/playback/docs/how-to-recreate-dashboard-queries-via-the-api-1#scale-factor)**

```yaml
queries:
# ...
- type: average
  metric: SCALE_FACTOR
```

## Building

### Coding Conventions

#### Style Guidelines

While not strictly enforced, the basic preferred editor settings are set in the
[.editorconfig](./.editorconfig). Other than this, no style guidelines are
currently imposed.

#### Static Analysis

This project uses both [`go vet`](https://pkg.go.dev/cmd/vet) and
[`staticcheck`](https://staticcheck.io/) to perform static code analysis. These
checks are run via [`precommit`](https://pre-commit.com) on all commits. Though
this can be bypassed on local commit, both tasks are also run during
[the `validate` workflow](./.github/workflows/validate.yml) and must have no
errors in order to be merged.

#### Commit Messages

Commit messages must follow [the conventional commit format](https://www.conventionalcommits.org/en/v1.0.0/).
Again, while this can be bypassed on local commit, it is strictly enforced in
[the `validate` workflow](./.github/workflows/validate.yml).

The basic commit message structure is as follows.

```
<type>[optional scope][!]: <description>

[optional body]

[optional footer(s)]
```

In addition to providing consistency, the commit message is used by
[svu](https://github.com/caarlos0/svu) during
[the release workflow](./.github/workflows/release.yml). The presence and values
of certain elements within the commit message affect auto-versioning. For
example, the `feat` type will bump the minor version. Therefore, it is important
to use the guidelines below and carefully consider the content of the commit
message.

Please use one of the types below.

- `feat` (bumps minor version)
- `fix` (bumps patch version)
- `chore`
- `build`
- `docs`
- `test`

Any type can be followed by the `!` character to indicate a breaking change.
Additionally, any commit that has the text `BREAKING CHANGE:` in the footer will
indicate a breaking change.

### Local Development

For local development, simply use `go build` and `go run`. For example,

```bash
go build cmd/bitmovin/bitmovin.go
```

Or

```bash
go run cmd/bitmovin/bitmovin.go
```

If you prefer, you can also use [`goreleaser`](https://goreleaser.com/) with
the `--single-target` option to build the binary for the local `GOOS` and
`GOARCH` only.

```bash
goreleaser build --single-target
```

### Releases

Releases are built and packaged using [`goreleaser`](https://goreleaser.com/).
By default, a new release will be built automatically on any push to the `main`
branch. For more details, review the [`.goreleaser.yaml`](./.goreleaser.yaml)
and [the `goreleaser` documentation](https://goreleaser.com/intro/).

The [svu](https://github.com/caarlos0/svu) utility is used to generate the next
tag value [based on commit messages](https://github.com/caarlos0/svu#commit-messages-vs-what-they-do).

### GitHub Workflows

This project utilizes GitHub workflows to perform actions in response to
certain GitHub events.

| Workflow | Events | Description
| --- | --- | --- |
| [validate](./.github/workflows/validate.yml) | `push` | Runs [precommit](https://pre-commit.com) to perform static analysis and runs [commitlint](https://commitlint.js.org/#/) to validate the last commit message |
| [build](./.github/workflows/build.yml) | `push`, `pull_request` | Builds and tests code |
| [release](./.github/workflows/release.yml) | `push` to `main` branch | Generates a new tag using [svu](https://github.com/caarlos0/svu) and runs [`goreleaser`](https://goreleaser.com/) |
| [repolinter](./.github/workflows/repolinter.yml) | `pull_request` | Enforces repository content guidelines |

## Support

New Relic has open-sourced this project. This project is provided AS-IS WITHOUT
WARRANTY OR DEDICATED SUPPORT. Issues and contributions should be reported to
the project here on GitHub.

We encourage you to bring your experiences and questions to the
[Explorers Hub](https://discuss.newrelic.com/) where our community members
collaborate on solutions and new ideas.

### Privacy

At New Relic we take your privacy and the security of your information
seriously, and are committed to protecting your information. We must emphasize
the importance of not sharing personal data in public forums, and ask all users
to scrub logs and diagnostic information for sensitive information, whether
personal, proprietary, or otherwise.

We define “Personal Data” as any information relating to an identified or
identifiable individual, including, for example, your name, phone number, post
code or zip code, Device ID, IP address, and email address.

For more information, review [New Relic’s General Data Privacy Notice](https://newrelic.com/termsandconditions/privacy).

### Contribute

We encourage your contributions to improve this project! Keep in mind that
when you submit your pull request, you'll need to sign the CLA via the
click-through using CLA-Assistant. You only have to sign the CLA one time per
project.

If you have any questions, or to execute our corporate CLA (which is required
if your contribution is on behalf of a company), drop us an email at
opensource@newrelic.com.

**A note about vulnerabilities**

As noted in our [security policy](../../security/policy), New Relic is committed
to the privacy and security of our customers and their data. We believe that
providing coordinated disclosure by security researchers and engaging with the
security community are important means to achieve our security goals.

If you believe you have found a security vulnerability in this project or any of
New Relic's products or websites, we welcome and greatly appreciate you
reporting it to New Relic through [HackerOne](https://hackerone.com/newrelic).

If you would like to contribute to this project, review [these guidelines](./CONTRIBUTING.md).

To all contributors, we thank you!  Without your contribution, this project
would not be what it is today.

### License

The New Relic Bitmovin Analytics project is licensed under the
[Apache 2.0](http://apache.org/licenses/LICENSE-2.0.txt) License.

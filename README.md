<a href="https://opensource.newrelic.com/oss-category/#community-project"><picture><source media="(prefers-color-scheme: dark)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/dark/Community_Project.png"><source media="(prefers-color-scheme: light)" srcset="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Project.png"><img alt="New Relic Open Source community project banner." src="https://github.com/newrelic/opensource-website/raw/main/src/images/categories/Community_Project.png"></picture></a>

# NRI Bitmovin Analytics Integration

![GitHub forks](https://img.shields.io/github/forks/newrelic/nri-bitmovin-analytics?style=social)
![GitHub stars](https://img.shields.io/github/stars/newrelic/nri-bitmovin-analytics?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/newrelic/nri-bitmovin-analytics?style=social)

![GitHub all releases](https://img.shields.io/github/downloads/newrelic/nri-bitmovin-analytics/total)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/newrelic/nri-bitmovin-analytics)
![GitHub last commit](https://img.shields.io/github/last-commit/newrelic/nri-bitmovin-analytics)
![GitHub Release Date](https://img.shields.io/github/release-date/newrelic/nri-bitmovin-analytics)

![GitHub issues](https://img.shields.io/github/issues/newrelic/nri-bitmovin-analytics)
![GitHub issues closed](https://img.shields.io/github/issues-closed/newrelic/nri-bitmovin-analytics)
![GitHub pull requests](https://img.shields.io/github/issues-pr/newrelic/nri-bitmovin-analytics)
![GitHub pull requests closed](https://img.shields.io/github/issues-pr-closed/newrelic/nri-bitmovin-analytics)

This integration uses Bitmovin Analytics API to pull in the below metrics and push them into New Relic.

Supported Metrics:
1. max_concurrent_viewers
2. avg_rebuffer_percentage
3. cnt_play_attempts
4. cnt_video_start_failures
5. avg_video_startup_time_ms
6. avg_video_bitrate_mbps


## Standalone

The Standalone environment runs the data pipelines as an independant service, either on-premises or cloud instances like AWS EC2. It can run on Linux, macOS, Windows, and any OS with support for GoLang.

### Prerequisites

- Go 1.20 or later.

### Build

Open a terminal, CD to `cmd/standalone`, and run:

```
$ go build
```

### Configuring the Pipeline

The standalone environment requires a YAML file for pipeline configuration. The required keys are:

- `interval`: Integer. Time in seconds between requests(should be same as the schedule / cron).
- `exporter`: `nrmetrics`.
- `bitmovin_api_key`: String. Bitmovin API Key
- `bitmovin_license_key`: String. Bitmovin License Key.
- `bitmovin_tenant_org`: String. Bitmovin Tenant Org.
- `nr_account_id`: String.
- `nr_api_key`: String. Api key for writing.
- `nr_endpoint`: String. New Relic endpoint region. Either `US` or `EU`. Optional, default value is `US`.

Check `config/example_config.yaml` for a configuration example.


### Running the Pipeline

Just run the following command from the build folder:

```
$ ./standalone path/to/config.yaml
```

To run the pipeline on system start, check your specific system init documentation.

## Lambda

The Lambda environment runs the data pipeline in AWS Lambda instances. It's located in the `lambda` folder, and is divided into 3 binaries: `lambda/receiver`, `lambda/processor` and `lambda/exporter`.

### Prerequisites

- An AWS account.
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) tool.
- Go 1.20 or later.
- GNU Make.

### Setting Up AWS

1. Create 3 lambdas for **Receiver**, **Processor** and **Exporter**. Runtime `Go 1.x`, arch `x86_64`, and handler names `receiver`, `processor` and `exporter`.
2. Create an SQS for **ProcessorToExporter**, type Standard, condition OnSuccess.
3. Open **Receiver** lambda config->permissions->execution role. Add another permission->create inline policy, add Lambda write permissions `InvokeAsync` and `InvokeFunction`.
4. Edit **Receiver** lambda config, add as a destination another lambda, the **Processor**, with async invocation.
5. Open **Processor** lambda config->permissions->execution role. Add another permission->create inline policy, add SQS write permissions.
6. Edit **Processor** lambda config, add as a destination the SQS **ProcessorToExporter**.
7. Open **Exporter** lambda config->permissions->execution role. Add another permission->create inline policy, add SQS read and write permissions.
8. Edit **Exporter** lambda config, add as a trigger the SQS **ProcessorToExporter**.

Note: when creating and configuring the SQS service and trigger, make sure to set the timing and batching options you will need. For example, a time interval of 5 minutes and batching of 50 events.

### Build & Deploy

Open a terminal, CD to `lambda`, and run:

```
$ make recv=RECEIVER proc=PROCESSOR expt=EXPORTER
```

Where *RECEIVER*, *PROCESSOR*, and *EXPORTER*, are the AWS Lambda functions you just created in the previous step.

### Configuring the Pipeline

A Lambda pipeline requieres some configuration keys to be set as **environment variables**. To set up environment variables, go to AWS console, Lambda->Functions, click your function, Configuration->Environment variables:

Environment Variables to be set on the Receiver function:

- `interval`: Integer. Time in seconds between requests(should be same as the schedule / cron).
- `exporter`: `nrmetrics`.
- `bitmovin_api_key`: String. Bitmovin API Key
- `bitmovin_license_key`: String. Bitmovin License Key.
- `bitmovin_tenant_org`: String. Bitmovin Tenant Org.

Environment Variables to be set on the Exporter function:

- `nr_account_id`: String. Account ID. Only requiered for `nrevents` and `nrapi` exporters.
- `nr_api_key`: String. Api key for writing.
- `nr_endpoint`: String. New Relic endpoint region. Either `US` or `EU`. Optional, default value is `US`.

### Running the Pipeline

Finally, to start running the pipeline you will need an EventBridge rule. Add a trigger for the **Receiver** lambda, select EventBridge as the source, create new rule, schedule expression `rate(1 minute)` (or the time you desire).

### Testing

Instead of running the pipeline with an EventBridge rule, you can just send async invocations to the **Receiver** lambda from the command line, using the following command:

```
$ aws lambda invoke-async --function-name RECEIVER --invoke-args INPUT.json
```

Where *RECEIVER* is the **Receiver** lambda name and *INPUT.json* is a file containing any JSON (the input event will be ignored by the receiver).

This will simulate a timer event and trigger the pipeline.

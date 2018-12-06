# Sidekiq Stats to Cloudwatch

## Description

This is a serverless function which polls a sidekiq queue API endpoint, transferring the data to AWS Cloudwatch Metrics.

## Requirements

1. Serverless Framework
1. Golang 1.x
1. A sidekiq stats endpoint with the following format. Set this to the environment variable `SIDEKIQ_STATS_URL`
    ```
    {
      "enqueued": 0,
      "busy": 0,
      "retries": 691,
      "queues": [
          {
              "name": "shipment_creation",
              "size": 0,
              "latency": 0
          }
      ]
    }
    ```

## Development Setup

Install serverless

```
npm install -g serverless
```

Install dependencies:

```
dep ensure
```

## Env Variables

 Name | Required | Default | Description |
|------|----------|---------|-------------|
| `SIDEKIQ_STATS_URL` | y | `-` | Sidekiq queue length endpoint to poll against  |
| `AWS_REGION` | y | `ap-southeast-1` | AWS Region where Cloudwatch Metrics are stored  |
| `APP_NAME` | n | `-` | Used to namespace Cloudwatch Metrics  |

## Deploy

Important: Make sure you have `production.yml` file present

Then:

```
make deploy
```

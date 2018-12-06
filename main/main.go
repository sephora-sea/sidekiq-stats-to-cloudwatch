package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/sephora-sea/sidekiq-stats-to-cloudwatch/config"
)

var myClient = &http.Client{Timeout: 10 * time.Second}
var metricTypes = []string{
	"enqueued",
	"busy",
	"retries",
}

// QueueLengthResponse defines the structure of the endpoint where this function will read the stats against
// Example:
//    {
//			"enqueued": 0,
//			"busy": 0,
//			"retries": 691,
//			"queues": [
//			    {
//			        "name": "shipment_creation",
//			        "size": 0,
//			        "latency": 0
//			    }
//			]
//		}
type QueueLengthResponse struct {
	Enqueued float64  `json:"enqueued"`
	Busy     float64  `json:"busy"`
	Retries  float64  `json:"retries"`
	Queues   []*Queue `json:"queues"`
}

type Queue struct {
	Name    string  `json:"name"`
	Size    float64 `json:"size"`
	Latency float64 `json:"latency"`
}

// Handler reads metrics from an endpoint (of your app), and puts them into cloudwatch metrics
func Handler(ctx context.Context) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetInstance().AWSRegion),
	})

	if err != nil {
		panic(err)
	}

	svc := cloudwatch.New(sess)

	response := QueueLengthResponse{}
	err = getResult(&response)
	if err != nil {
		panic(err)
	}

	log.Printf("%+v", response)

	for _, metricType := range metricTypes {
		output, err := putMetricData(svc, metricType, "Count", response.Enqueued)
		if err != nil {
			panic(err)
		}
		log.Printf("%+v", output)
	}

	for _, queue := range response.Queues {
		output, err := putDimensionedMetric(svc, "Queue Size", "Count", queue.Size, []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("Queue"),
				Value: aws.String(queue.Name),
			},
		})
		if err != nil {
			panic(err)
		}
		log.Printf("%+v", output)

		output, err = putDimensionedMetric(svc, "Queue Latency", "Seconds", queue.Latency, []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("Queue"),
				Value: aws.String(queue.Name),
			},
		})
		if err != nil {
			panic(err)
		}
		log.Printf("%+v", output)
	}
}

func putMetricData(cloudwatchSvc *cloudwatch.CloudWatch, attribute, unit string, value float64) (result *cloudwatch.PutMetricDataOutput, err error) {
	timeNow := time.Now()
	result, err = cloudwatchSvc.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(attribute),
				Timestamp:  &timeNow,
				Unit:       aws.String(unit),
				Value:      &value,
			},
		},
		Namespace: aws.String(config.GetInstance().AppName),
	})

	return
}

func putDimensionedMetric(cloudwatchSvc *cloudwatch.CloudWatch, attribute, unit string, value float64, dimensions []*cloudwatch.Dimension) (result *cloudwatch.PutMetricDataOutput, err error) {
	timeNow := time.Now()
	result, err = cloudwatchSvc.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String(attribute),
				Dimensions: dimensions,
				Timestamp:  &timeNow,
				Unit:       aws.String(unit),
				Value:      &value,
			},
		},
		Namespace: aws.String("OMS Sidekiq Queues"),
	})

	return
}

func getResult(target *QueueLengthResponse) (err error) {
	response, err := myClient.Get(config.GetInstance().SidekiqStatsURL)
	if err != nil {
		return
	}

	defer response.Body.Close()
	json.NewDecoder(response.Body).Decode(target)

	return
}

func main() {
	lambda.Start(Handler)
}

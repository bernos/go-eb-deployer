package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"log"
	"time"
)

type EBService struct {
	client *elasticbeanstalk.ElasticBeanstalk
}

func NewEBService(client *elasticbeanstalk.ElasticBeanstalk) *EBService {
	return &EBService{
		client: client,
	}
}

func (svc *EBService) GetEnvironments(applicationName string) ([]*elasticbeanstalk.EnvironmentDescription, error) {
	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		ApplicationName: aws.String(applicationName),
		IncludeDeleted:  aws.Boolean(false),
	}

	if resp, err := svc.client.DescribeEnvironments(params); err == nil {
		return resp.Environments, nil
	} else {
		return nil, err
	}
}

func (svc *EBService) ApplicationVersionExists(applicationName string, version string) (bool, error) {
	params := &elasticbeanstalk.DescribeApplicationVersionsInput{
		ApplicationName: aws.String(applicationName),
		VersionLabels: []*string{
			aws.String(version),
		},
	}

	if resp, err := svc.client.DescribeApplicationVersions(params); err == nil {
		for _, v := range resp.ApplicationVersions {
			if *v.VersionLabel == version {
				return true, nil
			}
		}
		return false, nil
	} else {
		return false, err
	}
}

func (svc *EBService) TerminateEnvironment(environmentId string) error {
	params := &elasticbeanstalk.TerminateEnvironmentInput{
		EnvironmentID: &environmentId,
	}

	if _, err := svc.client.TerminateEnvironment(params); err != nil {
		return err
	}
	// TODO: Wait for env to be terminated...
	return nil
}

func (svc *EBService) CreateApplicationVersion(applicationName string, version string, bucket string, key string) error {
	params := &elasticbeanstalk.CreateApplicationVersionInput{
		ApplicationName:       aws.String(applicationName), // Required
		VersionLabel:          aws.String(version),         // Required
		AutoCreateApplication: aws.Boolean(true),
		SourceBundle: &elasticbeanstalk.S3Location{
			S3Bucket: aws.String(bucket),
			S3Key:    aws.String(key),
		},
	}

	if _, err := svc.client.CreateApplicationVersion(params); err != nil {
		return err
	}

	return nil
}

func (svc *EBService) WaitForEnvironment(applicationName string, environmentName string, fnTest func(*elasticbeanstalk.EnvironmentDescription) bool) error {
	ready := false

	params := &elasticbeanstalk.DescribeEnvironmentsInput{
		ApplicationName: aws.String(applicationName),
		EnvironmentNames: []*string{
			aws.String(environmentName),
		},
	}

	for {
		if resp, err := svc.client.DescribeEnvironments(params); err == nil {
			if len(resp.Environments) == 0 {
				return fmt.Errorf("Could not wait for environment '%s' as it could not be found", environmentName)
			}
			ready = fnTest(resp.Environments[0])
		} else {
			return err
		}

		if ready {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func (svc *EBService) LogEnvironmentEvents(applicationName string, environmentName string, done <-chan struct{}) {

	messages := make(chan *elasticbeanstalk.EventDescription)
	since := time.Now().Add(time.Minute * -1)

	go func() {

		defer close(messages)

		for {
			params := &elasticbeanstalk.DescribeEventsInput{
				ApplicationName: aws.String(applicationName),
				EnvironmentName: aws.String(environmentName),
				Severity:        aws.String("TRACE"),
				StartTime:       aws.Time(since),
			}

			if resp, err := svc.client.DescribeEvents(params); err == nil {
				for i := len(resp.Events) - 1; i >= 0; i-- {

					select {
					case messages <- resp.Events[i]:
					case <-done:
						return
					}
					since = resp.Events[i].EventDate.Add(time.Second)
				}
			}

			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for msg := range messages {
			log.Println(msg.EventDate.String() + " - " + *msg.Message)
		}
	}()
}

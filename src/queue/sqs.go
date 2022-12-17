package queue

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/dstoy/tunnel/src/config"
)

type SQS struct {
	url *string
	svc *sqs.SQS
}

/**
* Initialize the SQS queue
 */
func (this *SQS) Connect(queueConfig *config.QueueConfig) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            this.getConfig(queueConfig),
		SharedConfigState: session.SharedConfigEnable,
	}))

	this.svc = sqs.New(sess)

	url, err := this.getQueueUrl(&queueConfig.Url)
	if err != nil {
		return err
	}

	this.url = url
	log.Println("Listening for events on queue:", *this.url)

	return nil
}

/**
 * Build the AWS credentials
 */
func (this *SQS) getConfig(queue *config.QueueConfig) aws.Config {
	var config = &aws.Config{}

	if queue.Region != "" {
		config = config.WithRegion(queue.Region)
	}

	if queue.KeyId != "" && queue.Secret != "" {
		creds := credentials.NewStaticCredentials(
			queue.KeyId,
			queue.Secret,
			"",
		)

		config = config.WithCredentials(creds)
	}

	return *config
}

/**
 * Get a message from the queue
 */
func (this *SQS) GetMessage() *Message {
	for {
		res, err := this.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            this.url,
			MaxNumberOfMessages: aws.Int64(1),
			WaitTimeSeconds:     aws.Int64(20),
			VisibilityTimeout:   aws.Int64(0),
		})
		if err != nil {
			log.Println("Error receiving SQS message: ", err)
			time.Sleep(10 * time.Second)
			continue
		}

		if len(res.Messages) > 0 {
			var raw = res.Messages[0]
			var message = &Message{
				Event: *raw.Body,
				Id:    *raw.ReceiptHandle,
			}

			return message
		}
	}
}

/**
 * Delete a message
 */
func (this *SQS) DeleteMessage(message *Message) error {
	_, err := this.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      this.url,
		ReceiptHandle: &message.Id,
	})
	return err
}

/**
 * Return the queue url
 */
func (this *SQS) getQueueUrl(queue *string) (*string, error) {
	url, err := this.svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: queue})
	if err != nil {
		// Create the queue if it does not exist
		var code = err.(awserr.Error).Code()
		if code == sqs.ErrCodeQueueDoesNotExist {
			log.Printf("Queue '%s' is missing. Creating it... \n", *queue)

			_, err := this.svc.CreateQueue(&sqs.CreateQueueInput{
				QueueName: queue,
				Attributes: map[string]*string{
					"MessageRetentionPeriod": aws.String("86400"),
				},
			})
			if err != nil {
				return nil, err
			}

			// Try to get the url again after it was created
			url, err = this.svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: queue})
			if err != nil {
				return nil, err
			}

			log.Printf("Queue '%s' successfully created \n", *queue)
		} else {
			return nil, err
		}
	}

	return url.QueueUrl, nil
}

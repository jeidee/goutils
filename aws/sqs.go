package aws

import (
	"errors"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	session *session.Session
	svc     *sqs.SQS
	qurl    string
}

var awssqs *SQS

func GetSQS(params ...string) (*SQS, error) {
	if len(params) == 0 {
		if awssqs == nil {
			return nil, errors.New("SQS is not initialized")
		}
		return awssqs, nil
	}

	region := ""
	credentialsProfile := ""
	qurl := ""

	if len(params) >= 1 {
		qurl = params[0]
	}

	if len(params) >= 2 {
		region = params[1]
	}

	if len(params) >= 3 {
		credentialsProfile = params[2]
	}

	awssqs = &SQS{}

	// Create a new session
	sess, err := func() (*session.Session, error) {
		if credentialsProfile != "" {
			return session.NewSession(&aws.Config{
				Region:      aws.String(region),
				Credentials: credentials.NewSharedCredentials("", credentialsProfile),
			})
		}

		return session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	}()

	if err != nil {
		return nil, err
	}

	awssqs.session = sess
	awssqs.svc = sqs.New(sess)
	awssqs.qurl = qurl

	return awssqs, nil
}

func (o *SQS) SendMessage(message string,
	attributes map[string]*sqs.MessageAttributeValue,
	delaySeconds int64) (*sqs.SendMessageOutput, error) {

	ret, err := o.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:       aws.String(message),
		QueueUrl:          aws.String(o.qurl),
		MessageAttributes: attributes,
		DelaySeconds:      aws.Int64(delaySeconds),
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *SQS) SendMessageBatch(entries []*sqs.SendMessageBatchRequestEntry) (*sqs.SendMessageBatchOutput, error) {

	ret, err := o.svc.SendMessageBatch(&sqs.SendMessageBatchInput{
		QueueUrl: aws.String(o.qurl),
		Entries:  entries,
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *SQS) ReceiveMessage(isAll bool,
	isApproximateFirstReceveiTimestamp bool,
	isApproximateReceiveCount bool,
	isSenderId bool,
	isSentTimestamp bool,
	maxNumberOfMessages int64,
	messageAttributesNames []*string,
	visibilityTimeout int64,
	waitTimeSeconds int64) (*sqs.ReceiveMessageOutput, error) {

	attributesNames := make([]*string, 5)
	if isAll {
		attributesNames = append(attributesNames, aws.String("ALL"))
	}
	if isApproximateFirstReceveiTimestamp {
		attributesNames = append(attributesNames, aws.String("ApproximateFirstReceiveTimestamp"))
	}
	if isApproximateReceiveCount {
		attributesNames = append(attributesNames, aws.String("ApproximateReceiveCount"))
	}
	if isSenderId {
		attributesNames = append(attributesNames, aws.String("SenderId"))
	}
	if isSentTimestamp {
		attributesNames = append(attributesNames, aws.String("SentTimestamp"))
	}

	ret, err := o.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames:        attributesNames,
		MaxNumberOfMessages:   aws.Int64(maxNumberOfMessages),
		MessageAttributeNames: messageAttributesNames,
		QueueUrl:              aws.String(o.qurl),
		VisibilityTimeout:     aws.Int64(visibilityTimeout),
		WaitTimeSeconds:       aws.Int64(waitTimeSeconds),
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *SQS) DeleteMessage(receiptHandle string) (*sqs.DeleteMessageOutput, error) {

	ret, err := o.svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(o.qurl),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (o *SQS) DeleteMessageBatch(entries []*sqs.DeleteMessageBatchRequestEntry) (*sqs.DeleteMessageBatchOutput, error) {

	ret, err := o.svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		QueueUrl: aws.String(o.qurl),
		Entries:  entries,
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func GetMessageAttributeValueFromInt(value int) *sqs.MessageAttributeValue {
	return &sqs.MessageAttributeValue{
		DataType:    aws.String("Number"),
		StringValue: aws.String(strconv.Itoa(value)),
	}
}

func GetMessageAttributeValueFromInt64(value int64) *sqs.MessageAttributeValue {
	return &sqs.MessageAttributeValue{
		DataType:    aws.String("Number"),
		StringValue: aws.String(strconv.FormatInt(value, 10)),
	}
}

func GetMessageAttributeValueFromString(value string) *sqs.MessageAttributeValue {
	return &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(value),
	}
}

func GetMessageBatchRequestEntry(
	message string,
	delaySeconds int64,
	id string,
	attributes map[string]*sqs.MessageAttributeValue) *sqs.SendMessageBatchRequestEntry {

	return &sqs.SendMessageBatchRequestEntry{
		MessageBody:       aws.String(message),
		DelaySeconds:      aws.Int64(delaySeconds),
		Id:                aws.String(id),
		MessageAttributes: attributes,
	}
}

func GetDeleteMessageBatchRequestEntry(
	id string,
	receiptHandle string) *sqs.DeleteMessageBatchRequestEntry {

	return &sqs.DeleteMessageBatchRequestEntry{
		Id:            aws.String(id),
		ReceiptHandle: aws.String(receiptHandle),
	}
}

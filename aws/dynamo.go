package aws

import (
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDB 다이나모 DB 관련 기능 포함
type DynamoDB struct {
	Session *session.Session
	DB      *dynamodb.DynamoDB
}

var dynamoDB *DynamoDB

func GetDynamoDB(params ...string) (*DynamoDB, error) {
	if len(params) == 0 {
		if dynamoDB == nil {
			return nil, errors.New("DB is not initialized")
		}
		return dynamoDB, nil
	}

	region := ""
	credentialsFile := ""
	credentialsProfile := ""

	if len(params) >= 1 {
		region = params[0]
	}

	if len(params) >= 2 {
		credTokens := strings.Split(params[1], ":")
		if len(credTokens) == 2 {
			credentialsFile = credTokens[0]
			credentialsProfile = credTokens[1]
		} else {
			credentialsFile = credTokens[0]
		}
	}

	dynamoDB = &DynamoDB{}

	// Create a new session
	sess, err := func() (*session.Session, error) {
		if credentialsProfile != "" {
			return session.NewSession(&aws.Config{
				Region:      aws.String(region),
				Credentials: credentials.NewSharedCredentials(credentialsFile, credentialsProfile),
			})
		}

		return session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	}()

	if err != nil {
		return nil, err
	}

	dynamoDB.Session = sess
	dynamoDB.DB = dynamodb.New(dynamoDB.Session)

	return dynamoDB, nil
}

// ExistsTable 함수는 테이블 존재 여부를 확인한다.
func (o *DynamoDB) ExistsTable(table string) (bool, error) {
	descParams := &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	}
	_, err := o.DB.DescribeTable(descParams)

	if err != nil {
		return false, err
	}
	return true, nil
}

func GetStringFromAttributeValue(o *dynamodb.AttributeValue) string {
	if o == nil || o.S == nil {
		return ""
	}

	return *o.S
}

func GetIntFromAttributeValue(o *dynamodb.AttributeValue) int {
	if o == nil || o.N == nil {
		return -1
	}

	n, err := strconv.Atoi(*o.N)
	if err != nil {
		return -1
	}

	return n
}

func GetInt64FromAttributeValue(o *dynamodb.AttributeValue) int64 {
	if o == nil || o.N == nil {
		return -1
	}

	n, err := strconv.ParseInt(*o.N, 10, 0)
	if err != nil {
		return -1
	}

	return n
}

func GetBoolFromAttributeValue(o *dynamodb.AttributeValue) bool {
	if o == nil || o.BOOL == nil {
		// 기본값은 false
		return false
	}

	return *o.BOOL
}

func GetAttributeValueFromString(value string) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{S: aws.String(value)}
}

func GetAttributeValueUpdateFromString(value string) *dynamodb.AttributeValueUpdate {
	return &dynamodb.AttributeValueUpdate{
		Action: aws.String("PUT"),
		Value:  GetAttributeValueFromString(value),
	}
}

func GetAttributeValueUpdateFromInt(value int) *dynamodb.AttributeValueUpdate {
	return &dynamodb.AttributeValueUpdate{
		Action: aws.String("PUT"),
		Value:  GetAttributeValueFromInt(value),
	}
}

func GetAttributeValueFromInt(value int) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(value))}
}

func GetAttributeValueFromInt64(value int64) *dynamodb.AttributeValue {
	return &dynamodb.AttributeValue{N: aws.String(strconv.FormatInt(value, 10))}
}

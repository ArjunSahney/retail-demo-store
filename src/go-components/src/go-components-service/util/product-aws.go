// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package util

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Pro_sess, Pro_err = session.NewSession(&aws.Config{})

// DynamoDB table names passed via environment
var DbTableProducts = os.Getenv("DDB_TABLE_PRODUCTS")
var DbTableCategories = os.Getenv("DDB_TABLE_CATEGORIES")

// Allow DDB endpoint to be overridden to support amazon/dynamodb-local
var DdbEndpointOverride = os.Getenv("DDB_ENDPOINT_OVERRIDE")
var RunningLocal bool

var DynamoClient *dynamodb.DynamoDB

// Initialize clients
func init() {
	if len(DdbEndpointOverride) > 0 {
		RunningLocal = true
		log.Println("Creating DDB client with endpoint override: ", DdbEndpointOverride)
		creds := credentials.NewStaticCredentials("does", "not", "matter")
		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      aws.String("us-east-1"),
			Endpoint:    aws.String(DdbEndpointOverride),
		}
		DynamoClient = dynamodb.New(Pro_sess, awsConfig)
	} else {
		log.Println("Checking DDDDBDBDB", DdbEndpointOverride)
		RunningLocal = false
		DynamoClient = dynamodb.New(Pro_sess)
	}
}

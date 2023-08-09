// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package util

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// var sess *session.Session
var Sess, Err = session.NewSession(&aws.Config{})

var Pinpoint_app_id = os.Getenv("PINPOINT_APP_ID")
var Pinpoint_client = pinpoint.New(Sess)
var Ssm_client = ssm.New(Sess)

// Connect Stuff
func init() {

}

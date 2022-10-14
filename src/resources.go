package sato

import (
	_ "embed" //required for embed"
)

//go:embed resources/aws_sns_topic.policy.template
var awsSNSTopic []byte

//go:embed resources/aws_iam_role.template
var awsIamRole []byte

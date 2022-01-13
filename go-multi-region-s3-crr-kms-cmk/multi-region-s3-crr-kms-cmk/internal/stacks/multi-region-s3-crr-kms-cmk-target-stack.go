package stacks

import (
	"multi-region-s3-crr-kms-cmk/pkg/target"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MultiRegionS3CrrKmsCmkTargetStack struct {
	Stack                       awscdk.Stack
	TargetBucket                awss3.Bucket
	TargetKeyIdSsmParameterName string
	Region                      *string
}

func NewMultiRegionS3CrrKmsCmkTargetStack(scope constructs.Construct, id *string, props *awscdk.StackProps) MultiRegionS3CrrKmsCmkTargetStack {
	stack := awscdk.NewStack(scope, id, props)

	myTargetConstruct := target.NewMultiRegionS3CrrKmsCmkTarget(stack, jsii.String("MyTarget"))

	return MultiRegionS3CrrKmsCmkTargetStack{
		Stack:                       stack,
		TargetBucket:                myTargetConstruct.TargetBucket,
		TargetKeyIdSsmParameterName: myTargetConstruct.TargetKeyIdSsmParameterName,
		Region:                      stack.Region(),
	}
}

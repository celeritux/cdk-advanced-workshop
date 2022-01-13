package stacks

import (
	"multi-region-s3-crr-kms-cmk/pkg/source"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MultiRegionS3CrrKmsCmkSourceStackProps struct {
	StackProps                  awscdk.StackProps
	TargetBucket                awss3.Bucket
	TargetKeyIdSsmParameterName string
	TargetRegion                *string
}

type MultiRegionS3CrrKmsCmkSourceStack struct {
	Stack awscdk.Stack
}

func NewMultiRegionS3CrrKmsCmkSourceStack(scope constructs.Construct, id *string, props *MultiRegionS3CrrKmsCmkSourceStackProps) MultiRegionS3CrrKmsCmkSourceStack {
	stack := awscdk.NewStack(scope, id, &props.StackProps)

	source.NewMultiRegionS3CrrKmsCmkSource(stack, jsii.String("MySource"), &source.MultiRegionS3CrrKmsCmkSourceProps{
		TargetBucket:                props.TargetBucket,
		TargetKeyIdSsmParameterName: props.TargetKeyIdSsmParameterName,
		TargetRegion:                props.TargetRegion,
	})

	return MultiRegionS3CrrKmsCmkSourceStack{
		Stack: stack,
	}
}

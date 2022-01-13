package pipeline

import (
	"multi-region-s3-crr-kms-cmk/internal/stacks"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

func NewMultiRegionS3CrrKmsCmkStage(scope constructs.Construct, id *string, props awscdk.StageProps) awscdk.Stage {
	stage := awscdk.NewStage(scope, id, &props)

	targetStack := stacks.NewMultiRegionS3CrrKmsCmkTargetStack(stage, jsii.String("MultiRegionS3CrrKmsCmkTarget"), &awscdk.StackProps{
		Env: &awscdk.Environment{
			Account: props.Env.Account,
			Region:  jsii.String("us-west-2"),
		},
	})

	sourceStack := stacks.NewMultiRegionS3CrrKmsCmkSourceStack(stage, jsii.String("MultiRegionS3CrrKmsCmkSource"), &stacks.MultiRegionS3CrrKmsCmkSourceStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: props.Env.Account,
				Region:  jsii.String("us-west-1"),
			},
		},
		TargetBucket:                targetStack.TargetBucket,
		TargetKeyIdSsmParameterName: targetStack.TargetKeyIdSsmParameterName,
		TargetRegion:                targetStack.Region,
	})

	sourceStack.Stack.AddDependency(targetStack.Stack, nil)

	return stage
}

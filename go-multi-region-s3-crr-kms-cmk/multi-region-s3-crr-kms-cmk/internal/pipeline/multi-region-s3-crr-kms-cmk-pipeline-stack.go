package pipeline

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodecommit"
	"github.com/aws/aws-cdk-go/awscdk/v2/pipelines"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MultiRegionS3CrrKmsCmkPipelineStackProps struct {
	StackProps awscdk.StackProps
}

func NewMultiRegionS3CrrKmsCmkPipelineStack(scope constructs.Construct, id *string, props *MultiRegionS3CrrKmsCmkPipelineStackProps) {
	stack := awscdk.NewStack(scope, id, &props.StackProps)

	//import codecommit repository
	repo := awscodecommit.Repository_FromRepositoryArn(stack, jsii.String("CodeCommitRepo"), jsii.String("arn:aws:codecommit:us-west-1:442530952648:cdk-pipelines-workshop"))

	pipeline := pipelines.NewCodePipeline(stack, jsii.String("Pipeline"), &pipelines.CodePipelineProps{
		Synth: pipelines.NewShellStep(jsii.String("Synth"), &pipelines.ShellStepProps{
			Input: pipelines.CodePipelineSource_CodeCommit(repo, jsii.String("main"), &pipelines.CodeCommitSourceOptions{}),
			Commands: jsii.Strings(
				"npm install -g aws-cdk && goenv global $GOLANG_16_VERSION",
				"cdk synth",
			),
		}),
	})

	// This is where we add the application stages
	pipeline.AddStage(NewMultiRegionS3CrrKmsCmkStage(stack, jsii.String("Pre-Prod"), awscdk.StageProps{
		Env: &awscdk.Environment{
			Account: props.StackProps.Env.Account,
			Region:  jsii.String("us-west-1"),
		},
	}), &pipelines.AddStageOpts{})

	pipeline.AddStage(NewMultiRegionS3CrrKmsCmkStage(stack, jsii.String("Prod"), awscdk.StageProps{
		Env: &awscdk.Environment{
			Account: props.StackProps.Env.Account,
			Region:  jsii.String("us-west-1"),
		},
	}), &pipelines.AddStageOpts{
		Pre: &[]pipelines.Step{
			pipelines.NewManualApprovalStep(jsii.String("PromoteToProd"), &pipelines.ManualApprovalStepProps{}),
		},
	})
}

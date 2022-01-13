package main

import (
	"multi-region-s3-crr-kms-cmk/internal/pipeline"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
)

func main() {
	app := awscdk.NewApp(nil)

	accountId := "442530952648"
	region := "us-west-1"

	pipeline.NewMultiRegionS3CrrKmsCmkPipelineStack(app, jsii.String("MultiRegionS3CrrKmsCmkPipelineStack"), &pipeline.MultiRegionS3CrrKmsCmkPipelineStackProps{
		StackProps: awscdk.StackProps{
			Env: &awscdk.Environment{
				Account: &accountId,
				Region:  &region,
			},
		},
	})

	app.Synth(nil)
}

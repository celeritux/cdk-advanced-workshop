package target

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	"fmt"
)

type MultiRegionS3CrrKmsCmkTarget struct {
	TargetBucket                awss3.Bucket
	TargetKeyIdSsmParameterName string
}

func NewMultiRegionS3CrrKmsCmkTarget(scope constructs.Construct, id *string) MultiRegionS3CrrKmsCmkTarget {
	construct := constructs.NewConstruct(scope, id)

	targetKmsKey := awskms.NewKey(construct, jsii.String("MyTargetKey"), &awskms.KeyProps{
		KeySpec:  awskms.KeySpec_SYMMETRIC_DEFAULT,
		KeyUsage: awskms.KeyUsage_ENCRYPT_DECRYPT,
		Enabled:  jsii.Bool(true),
	})

	targetBucket := awss3.NewBucket(construct, jsii.String("MyTargetBucket"), &awss3.BucketProps{
		BucketName:    awscdk.PhysicalName_GENERATE_IF_NEEDED(),
		Encryption:    awss3.BucketEncryption_KMS,
		EncryptionKey: targetKmsKey,
		Versioned:     jsii.Bool(true),
	})

	stack := awscdk.Stack_Of(construct)
	parameterName := fmt.Sprintf("%s.MyTargetKeyId", *stack.StackName())

	awsssm.NewStringParameter(construct, jsii.String("MyTargetKeyIdSSMParam"), &awsssm.StringParameterProps{
		StringValue:   targetKmsKey.KeyArn(),
		ParameterName: &parameterName,
		Description:   jsii.String("The KMS Key Id for the target stack"),
		DataType:      awsssm.ParameterDataType_TEXT,
		Tier:          awsssm.ParameterTier_STANDARD,
		Type:          awsssm.ParameterType_STRING,
	})

	return MultiRegionS3CrrKmsCmkTarget{
		TargetBucket:                targetBucket,
		TargetKeyIdSsmParameterName: parameterName,
	}
}

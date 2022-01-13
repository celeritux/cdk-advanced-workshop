package source

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/customresources"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MultiRegionS3CrrKmsCmkSourceProps struct {
	TargetBucket                awss3.Bucket
	TargetKeyIdSsmParameterName string
	TargetRegion                *string
}

func NewMultiRegionS3CrrKmsCmkSource(scope constructs.Construct, id *string, props *MultiRegionS3CrrKmsCmkSourceProps) {
	construct := constructs.NewConstruct(scope, id)

	sourceKmsKey := awskms.NewKey(construct, jsii.String("MySourceKey"), &awskms.KeyProps{
		KeySpec:  awskms.KeySpec_SYMMETRIC_DEFAULT,
		KeyUsage: awskms.KeyUsage_ENCRYPT_DECRYPT,
		Enabled:  jsii.Bool(true),
	})

	stack := awscdk.Stack_Of(construct)

	parameterArn := stack.FormatArn(&awscdk.ArnComponents{
		Account:      stack.Account(),
		Region:       props.TargetRegion,
		Resource:     jsii.String("parameter"),
		ResourceName: &props.TargetKeyIdSsmParameterName,
		Service:      jsii.String("ssm"),
	})

	targetKeyLookupCr := customresources.NewAwsCustomResource(construct, jsii.String("TargetKeyLookup"), &customresources.AwsCustomResourceProps{
		OnUpdate: &customresources.AwsSdkCall{
			Service: jsii.String("SSM"),
			Action:  jsii.String("getParameter"),
			Parameters: map[string]string{
				"Name": props.TargetKeyIdSsmParameterName,
			},
			Region:             props.TargetRegion,
			PhysicalResourceId: customresources.PhysicalResourceId_Of(jsii.String(time.Now().String())),
		},
		Policy: customresources.AwsCustomResourcePolicy_FromSdkCalls(&customresources.SdkCallsPolicyOptions{
			Resources: jsii.Strings(*parameterArn),
		}),
	})

	role := awsiam.NewRole(construct, jsii.String("MyCrrRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("s3.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		Path:      jsii.String("/service-role/"),
	})

	sourceCfnBucket := awss3.NewCfnBucket(construct, jsii.String("MySourceBucket"), &awss3.CfnBucketProps{
		BucketName: jsii.String(strings.ToLower(fmt.Sprintf("%s-%s-source", *stack.StackName(), *stack.Account()))),
		BucketEncryption: awss3.CfnBucket_BucketEncryptionProperty{
			ServerSideEncryptionConfiguration: []awss3.CfnBucket_ServerSideEncryptionRuleProperty{
				{
					ServerSideEncryptionByDefault: awss3.CfnBucket_ServerSideEncryptionByDefaultProperty{
						SseAlgorithm:   jsii.String("aws:kms"),
						KmsMasterKeyId: sourceKmsKey.KeyId(),
					},
				},
			},
		},
		VersioningConfiguration: awss3.CfnBucket_VersioningConfigurationProperty{
			Status: jsii.String("Enabled"),
		},
		ReplicationConfiguration: awss3.CfnBucket_ReplicationConfigurationProperty{
			Role: role.RoleArn(),
			Rules: []awss3.CfnBucket_ReplicationRuleProperty{
				{
					Destination: awss3.CfnBucket_ReplicationDestinationProperty{
						Bucket: props.TargetBucket.BucketArn(),
						EncryptionConfiguration: awss3.CfnBucket_EncryptionConfigurationProperty{
							ReplicaKmsKeyId: targetKeyLookupCr.GetResponseField(jsii.String("Parameter.Value")),
						},
					},
					SourceSelectionCriteria: awss3.CfnBucket_SourceSelectionCriteriaProperty{
						SseKmsEncryptedObjects: awss3.CfnBucket_SseKmsEncryptedObjectsProperty{
							Status: jsii.String("Enabled"),
						},
					},
					Status: jsii.String("Enabled"),
				},
			},
		},
	})
	sourceBucket := awss3.Bucket_FromBucketName(construct, jsii.String("MyImportedSourceBucket"), sourceCfnBucket.BucketName())

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: jsii.Strings(*sourceBucket.BucketArn()),
		Actions:   jsii.Strings("s3:GetReplicationConfiguration", "s3:ListBucket"),
	}))

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: jsii.Strings(*sourceBucket.ArnForObjects(jsii.String("*"))),
		Actions:   jsii.Strings("s3:GetObjectVersionForReplication", "s3:GetObjectVersionAcl", "s3:GetObjectVersionTagging"),
	}))

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: jsii.Strings(*props.TargetBucket.ArnForObjects(jsii.String("*"))),
		Actions:   jsii.Strings("s3:ReplicateObject", "s3:ReplicateDelete", "s3:ReplicateTags"),
	}))

	//sourceKmsKey.GrantDecrypt(role)
	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: jsii.Strings(*sourceKmsKey.KeyArn()),
		Actions:   jsii.Strings("kms:Descrypt"),
	}))

	role.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: jsii.Strings(*targetKeyLookupCr.GetResponseField(jsii.String("Parameter.Value"))),
		Actions:   jsii.Strings("kms:Encrypt"),
	}))

}

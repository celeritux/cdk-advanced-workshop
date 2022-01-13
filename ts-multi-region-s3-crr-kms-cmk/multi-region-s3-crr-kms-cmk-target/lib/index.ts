import { Construct } from 'constructs';
import * as cdk from 'aws-cdk-lib';
import * as kms from 'aws-cdk-lib/aws-kms';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as ssm from 'aws-cdk-lib/aws-ssm';

export class MultiRegionS3CrrKmsCmkTarget extends Construct {
  public readonly targetBucket: s3.Bucket;
  public readonly targetKeyIdSsmParameterName: string;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    const targetKmsKey = new kms.Key(this, 'MyTargetKey', {
      keySpec: kms.KeySpec.SYMMETRIC_DEFAULT,
      keyUsage: kms.KeyUsage.ENCRYPT_DECRYPT,
      enabled: true
    });

    const targetBucket = new s3.Bucket(this, 'MyTargetBucket', {
      bucketName: cdk.PhysicalName.GENERATE_IF_NEEDED,
      encryption: s3.BucketEncryption.KMS,
      encryptionKey: targetKmsKey,
      versioned: true
    });

    //exporting Kms key Arn to be used by other stack/region
    const stack = cdk.Stack.of(this);
    const parameterName = `${stack.stackName}.MyTargetKeyId`;
    new ssm.StringParameter(this, 'MyTargetKeyIdSSMParam', {
      stringValue: targetKmsKey.keyArn,
      parameterName: parameterName,
      description: 'The KMS Key Id for the target stack',
      dataType: ssm.ParameterDataType.TEXT,
      tier: ssm.ParameterTier.STANDARD,
      type: ssm.ParameterType.STRING,
    });

    this.targetBucket = targetBucket;
    // We use the SSM parameter name, with the KMS key Arn, to lookup the KMS key value in other cross region stacks
    this.targetKeyIdSsmParameterName = parameterName;
  }
}

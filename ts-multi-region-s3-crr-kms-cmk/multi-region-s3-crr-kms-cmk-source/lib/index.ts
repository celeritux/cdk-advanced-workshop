import { Construct } from 'constructs';
import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3'
import * as kms from 'aws-cdk-lib/aws-kms';
import * as cr from 'aws-cdk-lib/custom-resources';
import * as iam from 'aws-cdk-lib/aws-iam';

export interface MultiRegionS3CrrKmsCmkSourceProps {
  readonly targetBucket: s3.Bucket;
  readonly targetKeyIdSsmParameterName: string;
  readonly targetRegion: string;
}

export class MultiRegionS3CrrKmsCmkSource extends Construct {

  constructor(scope: Construct, id: string, props: MultiRegionS3CrrKmsCmkSourceProps) {
    super(scope, id);

    const sourceKmsKey = new kms.Key(this, 'MySourceKey', {
      keySpec: kms.KeySpec.SYMMETRIC_DEFAULT,
      keyUsage: kms.KeyUsage.ENCRYPT_DECRYPT,
      enabled: true
    });

    const sourceBucket = new s3.Bucket(this, 'MySourceBucket', {
      bucketName: cdk.PhysicalName.GENERATE_IF_NEEDED,
      encryption: s3.BucketEncryption.KMS,
      encryptionKey: sourceKmsKey,
      versioned: true
    });

    const stack = cdk.Stack.of(this);
    const parameterArn = stack.formatArn({
        account: stack.account,
        region: props.targetRegion,
        resource: 'parameter',
        resourceName: props.targetKeyIdSsmParameterName,
        service: 'ssm'
    });

    const targetKeyLookupCR = new cr.AwsCustomResource(this, 'TargetKeyLookup', {
      onUpdate: {   // will also be called for a CREATE event
        service: 'SSM',
        action: 'getParameter',
        parameters: {
          Name: props.targetKeyIdSsmParameterName
        },
        region: props.targetRegion,
        physicalResourceId: cr.PhysicalResourceId.of(Date.now().toString())
      },
      policy: cr.AwsCustomResourcePolicy.fromSdkCalls({resources: [parameterArn]})
    });

    const role = new iam.Role(this, 'MyCrrRole', {
      assumedBy: new iam.ServicePrincipal('s3.amazonaws.com'),
      path: '/service-role/'
    });

    role.addToPolicy(new iam.PolicyStatement({
      resources: [sourceBucket.bucketArn],
      actions: ['s3:GetReplicationConfiguration', 's3:ListBucket'] }));

    role.addToPolicy(new iam.PolicyStatement({
      resources: [sourceBucket.arnForObjects('*')],
      actions: ['s3:GetObjectVersion', 's3:GetObjectVersionAcl', 's3:GetObjectVersionForReplication', 's3:GetObjectLegalHold', 's3:GetObjectVersionTagging', 's3:GetObjectRetention'] }));

    role.addToPolicy(new iam.PolicyStatement({
      resources: [props.targetBucket.arnForObjects('*')],
      actions: ['s3:ReplicateObject', 's3:ReplicateDelete', 's3:ReplicateTags', 's3:GetObjectVersionTagging'] }));

    role.addToPolicy(new iam.PolicyStatement({
      resources: [sourceKmsKey.keyArn],
      actions: ['kms:Decrypt'] }));
    //sourceKmsKey.grantDecrypt(role);

    role.addToPolicy(new iam.PolicyStatement({
      resources: [targetKeyLookupCR.getResponseField('Parameter.Value')],
      actions: ['kms:Encrypt']
    }));

    const replicaConfigurationProperty: s3.CfnBucket.ReplicationConfigurationProperty = {
      role: role.roleArn,
      rules: [{
        destination: {
          bucket: props.targetBucket.bucketArn,
          encryptionConfiguration: {
            replicaKmsKeyId: targetKeyLookupCR.getResponseField('Parameter.Value')
          }
        },
        sourceSelectionCriteria: {
          sseKmsEncryptedObjects: {
            status: 'Enabled'
          }
        },
        status: 'Enabled'
      }]
    };

    // Get the AWS CloudFormation resource for the source bucket and set replica configuration to the source bucket
    (sourceBucket.node.defaultChild as s3.CfnBucket).replicationConfiguration = replicaConfigurationProperty;
  }
}

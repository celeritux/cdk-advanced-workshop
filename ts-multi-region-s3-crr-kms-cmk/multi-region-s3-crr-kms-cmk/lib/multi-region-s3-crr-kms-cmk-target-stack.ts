import { Construct } from 'constructs';
import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { MultiRegionS3CrrKmsCmkTarget } from 'multi-region-s3-crr-kms-cmk-target';

export class MultiRegionS3CrrKmsCmkTargetStack extends cdk.Stack {
    public targetBucket: s3.Bucket;
    public targetKeyIdSsmParameterName: string;

    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const myTargetConstruct = new MultiRegionS3CrrKmsCmkTarget(this, 'MyTarget');
        this.targetBucket = myTargetConstruct.targetBucket;
        this.targetKeyIdSsmParameterName = myTargetConstruct.targetKeyIdSsmParameterName;
    }
}
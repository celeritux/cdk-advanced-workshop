import { Construct } from 'constructs';
import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import { MultiRegionS3CrrKmsCmkSource } from 'multi-region-s3-crr-kms-cmk-source';

interface S3StaticMultiRegionSourceStackProps extends cdk.StackProps {
    targetBucket: s3.Bucket,
    targetKeyIdSsmParameterName: string,
    targetRegion: string
}

export class MultiRegionS3CrrKmsCmkSourceStack extends cdk.Stack {
    public targetBucket: s3.Bucket;
    public targetKeyIdSsmParameterName: string;

    constructor(scope: Construct, id: string, props: S3StaticMultiRegionSourceStackProps) {
        super(scope, id, props);

        const mySourceConstruct = new MultiRegionS3CrrKmsCmkSource(this, 'MySource', {
            targetBucket: props.targetBucket,
            targetKeyIdSsmParameterName: props.targetKeyIdSsmParameterName,
            targetRegion: props.targetRegion
        });
    }
}
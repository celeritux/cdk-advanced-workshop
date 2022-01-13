import { Construct } from 'constructs';
import {Stage, StageProps} from 'aws-cdk-lib';
import { MultiRegionS3CrrKmsCmkTargetStack } from './multi-region-s3-crr-kms-cmk-target-stack';
import { MultiRegionS3CrrKmsCmkSourceStack } from './multi-region-s3-crr-kms-cmk-source-stack';

export class MultiRegionS3CrrKmsCmkStage extends Stage {
    constructor(scope: Construct, id: string, props?: StageProps) {
        super(scope, id, props);

        const targetStack = new MultiRegionS3CrrKmsCmkTargetStack(this, 'MultiRegionS3CrrKmsCmkTarget', {
            env: { account: props?.env?.account, region: "us-west-2"}
        });

        const sourceStack = new MultiRegionS3CrrKmsCmkSourceStack(this, 'MultiRegionS3CrrKmsCmkSource', {
            env: { account: props?.env?.account, region: 'us-west-1' },
            targetBucket: targetStack.targetBucket,
            targetKeyIdSsmParameterName: targetStack.targetKeyIdSsmParameterName,
            targetRegion: targetStack.region
        });

        sourceStack.addDependency(targetStack);
    }
}
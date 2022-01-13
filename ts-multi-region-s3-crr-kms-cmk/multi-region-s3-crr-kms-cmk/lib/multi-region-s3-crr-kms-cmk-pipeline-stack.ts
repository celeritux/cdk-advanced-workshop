import { Construct } from 'constructs';
import { Stack, StackProps, Stage } from 'aws-cdk-lib';
import * as codecommit from 'aws-cdk-lib/aws-codecommit';
import * as pipelines from 'aws-cdk-lib/pipelines';
import { MultiRegionS3CrrKmsCmkStage } from './multi-region-s3-crr-kms-cmk-stage';

export class MultiRegionS3CrrKmsCmkPipelineStack extends Stack {
    constructor(scope: Construct, id: string, props?: StackProps) {
        super(scope, id, props);

        const repo = codecommit.Repository.fromRepositoryName(this, 'CodeCommitRepo', 'cdk-pipelines-workshop');

        const pipeline = new pipelines.CodePipeline(this, 'Pipeline', {
            pipelineName: 'MultiRegionS3CrrKmsCmkPipeline',
            
            //how we build and synthesized the application
            synth: new pipelines.CodeBuildStep('SynthStep', {
                input: pipelines.CodePipelineSource.codeCommit(repo, 'main'),
                installCommands: [
                    'npm install -g aws-cdk'
                ],
                commands: [
                    'npm ci',
                    'npm run build',
                    'npx cdk synth'
                ]
            })
        });

        pipeline.addStage(new MultiRegionS3CrrKmsCmkStage(this, 'PreProd', {
            env: {account: props?.env?.account, region: 'us-west-2'}
        }));

        pipeline.addStage(new MultiRegionS3CrrKmsCmkStage(this, 'Prod', {
            env: {account: props?.env?.account, region: 'us-west-2'}
        }), {
            pre: [new pipelines.ManualApprovalStep('PromoteToProd')]
        });
    }
}
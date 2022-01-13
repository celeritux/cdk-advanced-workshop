#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { MultiRegionS3CrrKmsCmkPipelineStack } from '../lib/multi-region-s3-crr-kms-cmk-pipeline-stack';

const app = new cdk.App();
const accountId = '442530952648'

/*
export ACCOUNT_ID=$(aws sts get-caller-identity --output text --profile ww.dev.fvigato --query Account)
npx cdk bootstrap \
  --cloudformation-execution-policies arn:aws:iam::aws:policy/AdministratorAccess \
  aws://$ACCOUNT_ID/us-west-1 \
  aws://$ACCOUNT_ID/us-west-2
*/

new MultiRegionS3CrrKmsCmkPipelineStack(app, 'MultiRegionS3CrrKmsCmkPipelineStack', {
    env: {account: accountId, region: "us-west-2"}
})
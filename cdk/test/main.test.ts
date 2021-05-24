import '@aws-cdk/assert/jest';
import { App } from '@aws-cdk/core';
import { KinesisTestStack } from '../src/main';

test('Snapshot', () => {
  const app = new App();
  const stack = new KinesisTestStack(app, 'test');

  expect(stack).toHaveResource('AWS::S3::Bucket');
  expect(app.synth().getStackArtifact(stack.artifactId).template).toMatchSnapshot();
});
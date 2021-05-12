import * as glue from '@aws-cdk/aws-glue';
import * as iam from '@aws-cdk/aws-iam';
import * as kinesis from '@aws-cdk/aws-kinesis';
import * as firehose from '@aws-cdk/aws-kinesisfirehose';
import * as logs from '@aws-cdk/aws-logs';
import * as s3 from '@aws-cdk/aws-s3';
import * as cdk from '@aws-cdk/core';

function String(name: string): glue.Column {
    return {
        name: name,
        type: glue.Schema.STRING,
    };
}

function Int(name: string): glue.Column {
    return {
        name: name,
        type: glue.Schema.INTEGER,
    };
}

export class KinesisTestStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props: cdk.StackProps = {}) {
    super(scope, id, props);

    const bucket = new s3.Bucket(this, 'FirehoseDeliveryBucket',  {
      autoDeleteObjects: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      enforceSSL: true,
      encryption: s3.BucketEncryption.S3_MANAGED,
    });

    const glueDb = new glue.Database(this, 'KinesisTestGlueDatabase', {
      databaseName: 'kinesis_test',
    });

    const tableName = 'test_glue_table';

    new glue.Table(this, 'ApplogsGlueTable', {
            bucket: bucket,
            s3Prefix: 'kinesis_delivered/',
            database: glueDb,
            tableName: tableName,
            dataFormat: glue.DataFormat.PARQUET,
            columns: [
                String('foo'),
                Int('n'),
                Int('baz'), // this one we will generate data conversion errors with
            ],
            partitionKeys: [String('ingesteddate')],
        });

    const deliveryRole = new iam.Role(this, 'DeliveryRole', {
            assumedBy: new iam.ServicePrincipal('firehose.amazonaws.com'),
            externalIds: [cdk.Aws.ACCOUNT_ID],
        });

    const source = new kinesis.Stream(this, 'IngestStream', {
            streamName: 'kinesis-test-ingest-stream',
            shardCount: 1,
            retentionPeriod: cdk.Duration.days(1),
            encryption: kinesis.StreamEncryption.MANAGED,
        });

    source.grantRead(deliveryRole);

    deliveryRole.addToPolicy(
            new iam.PolicyStatement({
                resources: [source.streamArn],
                actions: [
                    'kinesis:DescribeStream',
                    'kinesis:GetShardIterator',
                    'kinesis:GetRecords',
                    'kinesis:ListShards',
                    'firehose:DeleteDeliveryStream',
                    'firehose:PutRecord',
                    'firehose:PutRecordBatch',
                    'firehose:UpdateDestination',
                ],
            })
        );

    deliveryRole.addToPolicy(
        new iam.PolicyStatement({
            resources: ['*'],
            actions: ['glue:GetTable', 'glue:GetTableVersion', 'glue:GetTableVersions'],
        })
    );

    deliveryRole.addToPolicy(
        new iam.PolicyStatement({
            resources: [`${bucket.bucketArn}/*`, bucket.bucketArn],
            actions: [
                's3:AbortMultipartUpload',
                's3:GetBucketLocation',
                's3:GetObject',
                's3:ListBucket',
                's3:ListBucketMultipartUploads',
                's3:PutObject',
            ],
        })
    );

    const logGroupName = `/aws/kinesisfirehose/kdf-${tableName}`;

    deliveryRole.addToPolicy(
        new iam.PolicyStatement({
            resources: [
                `arn:aws:logs:${cdk.Aws.REGION}:${cdk.Aws.ACCOUNT_ID}:log-group:${logGroupName}:log-stream:*`,
            ],
            actions: ['logs:PutLogEvents'],
        })
    );

    const logGroup = new logs.LogGroup(this, 'FirehoseLogGroup', {
        retention: logs.RetentionDays.ONE_WEEK,
        logGroupName: logGroupName,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const logStream = new logs.LogStream(this, 'FirehoseLogStream', {
        logGroup: logGroup,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const backupsLogStream = new logs.LogStream(this, 'BackupsLogStream', {
        logGroup: logGroup,
        removalPolicy: cdk.RemovalPolicy.DESTROY,
    });

    const delivery = new firehose.CfnDeliveryStream(this, 'FirehoseDelivery', {
            deliveryStreamType: 'KinesisStreamAsSource',
            kinesisStreamSourceConfiguration: {
                roleArn: deliveryRole.roleArn,
                kinesisStreamArn: source.streamArn,
            },
            extendedS3DestinationConfiguration: {
                bucketArn: bucket.bucketArn,
                roleArn: deliveryRole.roleArn,
                s3BackupMode: 'Enabled',
                s3BackupConfiguration: {
                    bucketArn: bucket.bucketArn,
                    roleArn: deliveryRole.roleArn,
                    prefix: 'raw_original_records/IngestedDate=!{timestamp:yyyy-MM-dd}/',
                    errorOutputPrefix:
                        'raw_backups_errors/!{firehose:error-output-type}/!{timestamp:yyyy-MM-dd}/',
                    cloudWatchLoggingOptions: {
                        enabled: true,
                        logGroupName:  logGroup.logGroupName,
                        logStreamName: backupsLogStream.logStreamName,
                    },
                },
                bufferingHints: {
                    intervalInSeconds: 60,
                },
                cloudWatchLoggingOptions: {
                    enabled: true,
                    logGroupName: logGroup.logGroupName,
                    logStreamName: logStream.logStreamName,
                },
                encryptionConfiguration: {
                    noEncryptionConfig: 'NoEncryption',
                },
                errorOutputPrefix:
                    'raw_firehose_errors/!{firehose:error-output-type}/!{timestamp:yyyy-MM-dd}/',
                prefix: 'kinesis_delivered/IngestedDate=!{timestamp:yyyy-MM-dd}/!{firehose:random-string}/',
                dataFormatConversionConfiguration: {
                    enabled: true,
                    inputFormatConfiguration: {
                        deserializer: {
                            openXJsonSerDe: {
                                caseInsensitive: true,
                            },
                        },
                    },
                    outputFormatConfiguration: {
                        serializer: {
                            parquetSerDe: {
                                compression: 'SNAPPY',
                                writerVersion: 'V1',
                            },
                        },
                    },
                    schemaConfiguration: {
                        catalogId: cdk.Aws.ACCOUNT_ID,
                        databaseName: glueDb.databaseName,
                        tableName: tableName,
                        versionId: 'LATEST',
                        roleArn: deliveryRole.roleArn,
                    },
                },
            },
        });
        delivery.node.addDependency(deliveryRole);

        new cdk.CfnOutput(this, 'KinesisTestingBucket', { value: bucket.bucketName });
  }
}

// for development, use account/region from cdk cli
const devEnv = {
  account: process.env.CDK_DEFAULT_ACCOUNT,
  region: process.env.CDK_DEFAULT_REGION,
};

const app = new cdk.App();

new KinesisTestStack(app, 'kinesis-test-stack-dev', { env: devEnv });
// new MyStack(app, 'my-stack-prod', { env: prodEnv });

app.synth();
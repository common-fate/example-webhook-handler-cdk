import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
// import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as go from "@aws-cdk/aws-lambda-go-alpha";

export class ExampleWebhookHandlerStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const lambda = new go.GoFunction(this, "handler", {
      entry: "cmd/lambda",
    });

    const functionUrl = lambda.addFunctionUrl({
      authType: cdk.aws_lambda.FunctionUrlAuthType.NONE,
    });

    new cdk.CfnOutput(this, "FunctionUrl", {
      value: functionUrl.url,
      description: "URL for the Lambda function",
    });
  }
}

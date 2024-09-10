# example-webhook-handler-cdk

An example webhook handler for Common Fate, deployed using the AWS CDK. It uses a preshared authorization header for authentication.

## How to use

1. Replace the authorization header and the permission set identifier in `cmd/lambda/main.go`.

2. Deploy the Lambda function: `npm run cdk deploy`.

3. Add a webhook handler to your Common Fate configuration, replacing the `url` and `headers` values with ones matching the Lambda function you have deployed:

  ```hcl
  resource "commonfate_webhook_integration" "example_lambda" {
    name = "Example Lambda"
    url  = "https://abcdef.lambda-url.ap-southeast-2.on.aws" // REPLACE THIS
    headers = [
      {
        key   = "Authorization",
        value = "abcdef" // REPLACE THIS
      },
    ]
    send_audit_log_events = true
  }
  ```


## Useful commands

* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `npx cdk deploy`  deploy this stack to your default AWS account/region
* `npx cdk diff`    compare deployed stack with current state
* `npx cdk synth`   emits the synthesized CloudFormation template

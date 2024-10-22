#!/bin/bash

source $(dirname "$0")/init.sh

AWS_STACK_NAME=${AWS_STACK_NAME:-"$APP_DIR_NAME"}
AWS_CF_DEPLOY_OPTS=${AWS_CF_DEPLOY_OPTS:-""}
AWS_S3_BUCKET_KEY=${AWS_S3_BUCKET_KEY:-"newrelic-bitmovin-lambda.zip"}

if [ ! -n "$AWS_S3_BUCKET_NAME" -o ! -n "$AWS_S3_BUCKET_KEY" ]; then
    err "must specify S3 bucket name and key for uploading lambda zip"
    exit 1
fi

println "\n%s" "-- DEPLOY ----------------------------------------------------------------------"
println "Root directory:                          $ROOT_DIR"
println "Dist directory:                          $DIST_DIR"
println "Stage directory:                         $STAGE_DIR"
println "AWS region:                              $AWS_REGION"
println "Stack Name:                              $AWS_STACK_NAME"
println "Deploy stack options:                    $AWS_CF_DEPLOY_OPTS"
println "S3 bucket name:                          $AWS_S3_BUCKET_NAME"
println "S3 bucket key:                           $AWS_S3_BUCKET_KEY"
println "%s\n" "--------------------------------------------------------------------------------"

println "Building lambda zip package..."
build_zip

println "Uploading lambda zip package..."
upload_zip

println "Deploying stack $AWS_STACK_NAME..."
aws cloudformation deploy \
    --stack-name $AWS_STACK_NAME \
    --template-file $ROOT_DIR/deployments/lambda/cf-template.yaml \
    --output table \
    --no-cli-pager \
    --color on \
    --parameter-overrides file://$ROOT_DIR/deployments/lambda/cf-params.json \
    --capabilities CAPABILITY_NAMED_IAM \
    $AWS_CF_DEPLOY_OPTS

println "Done."

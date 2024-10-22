#!/bin/bash

source $(dirname "$0")/init.sh

FUNCTION_NAME=${FUNCTION_NAME:-"NewRelicBitmovinAnalytics"}
AWS_LAMBDA_UPDATE_OPTS=${AWS_LAMBDA_UPDATE_OPTS:-""}
AWS_S3_BUCKET_KEY=${AWS_S3_BUCKET_KEY:-"newrelic-bitmovin-lambda.zip"}

if [ ! -n "$AWS_S3_BUCKET_NAME" -o ! -n "$AWS_S3_BUCKET_KEY" ]; then
    err "must specify S3 bucket name and key for uploading lambda zip"
    exit 1
fi

println "\n%s" "-- UPDATE ----------------------------------------------------------------------"
println "Root directory:                          $ROOT_DIR"
println "Dist directory:                          $DIST_DIR"
println "Stage directory:                         $STAGE_DIR"
println "AWS region:                              $AWS_REGION"
println "Function name:                           $FUNCTION_NAME"
println "Update function options:                 $AWS_LAMBDA_UPDATE_OPTS"
println "S3 bucket name:                          $AWS_S3_BUCKET_NAME"
println "S3 bucket key:                           $AWS_S3_BUCKET_KEY"
println "%s\n" "--------------------------------------------------------------------------------"

println "Building lambda zip package..."
build_zip

println "Uploading lambda zip package..."
upload_zip

aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --s3-bucket $AWS_S3_BUCKET_NAME \
    --s3-key $AWS_S3_BUCKET_KEY \
    --output table \
    --no-cli-pager \
    --color on \
    $AWS_LAMBDA_UPDATE_OPTS

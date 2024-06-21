#!/bin/bash

if [ "$SCRIPT_DIR" = "" ]; then
  SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
fi

ROOT_DIR=$(dirname $(dirname ${SCRIPT_DIR}))
APP_DIR_NAME=$(basename "$ROOT_DIR")
DIST_DIR=$ROOT_DIR/dist
STAGE_DIR=$DIST_DIR/stage

function println {
  printf "$1\n" "$2"
}

function err {
  printf "\e[31m%s\e[0m\n\n" "$1"
  exit 1
}

function build_zip {
  CHECK=$(basename $(dirname $(dirname $STAGE_DIR)))

  if [ ! "$CHECK" == "newrelic-bitmovin-analytics" ]; then
    err "refusing to recursively remove $STAGE_DIR"
    exit 1
  fi

  rm -rf $STAGE_DIR
  mkdir -p $STAGE_DIR
  cp $ROOT_DIR/bin/bootstrap $STAGE_DIR
  cp -R $ROOT_DIR/configs $STAGE_DIR/configs
  cd $STAGE_DIR && zip -r $DIST_DIR/newrelic-bitmovin-lambda.zip .
}

function upload_zip {
  if [ ! -f "$DIST_DIR/newrelic-bitmovin-lambda.zip" ]; then
    err "can't find $DIST_DIR/newrelic-bitmovin-lambda.zip"
    exit 1
  fi

  aws s3 cp $DIST_DIR/newrelic-bitmovin-lambda.zip s3://$AWS_S3_BUCKET_NAME/$AWS_S3_BUCKET_KEY
}

if [ -z "$AWS_REGION" ]; then
  AWS_REGION=$(aws configure get region)
fi

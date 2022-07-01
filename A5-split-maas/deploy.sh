#!/bin/bash

while [[ $# -gt 0 ]]
do
    key="${1}"
    case ${key} in
    -t|--tracing)
        ENABLE_TRACING="${2}"
	echo "Tracing enabled: "$ENABLE_TRACING
        shift # past argument
        shift # past value
	;;
    *)    # unknown option
        shift # past argument
        ;;
    esac
    shift
done

sam package --template-file template.yaml --s3-bucket thesis-code-cloudformation --output-template-file packaged.yaml

aws cloudformation deploy \
        --template-file ./packaged.yaml \
        --stack-name test-stack \
	      --parameter-overrides EnableTracingParameter=$ENABLE_TRACING \
        --capabilities CAPABILITY_IAM

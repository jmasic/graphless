#!/bin/bash

FUNCTION=MainFunction
REGION=us-east-2

while [[ $# -gt 0 ]]
do
    key="${1}"
    case ${key} in
    -f|--function)
        FUNCTION="${2}"
        shift # past argument
        shift # past value
        ;;
    -r|--region)
        REGION="${2}"
        shift # past argument
        shift # past value
	    ;;
    -p|--payload)
        PAYLOAD_FILE="${2}"
        shift # past argument
        shift # past value
        ;;
    *)  # unknown option
        shift # past argument
        ;;
    esac
done

if [ -z $PAYLOAD_FILE ];then
    echo "Missing mandatory configuration file: --payload <payload-file>" >&2
    exit 1
fi

PAYLOAD_STRING=$(cat $PAYLOAD_FILE)
echo "Invoking Graphless platform with payload file $PAYLOAD_STRING"

aws lambda invoke --invocation-type Event \
		  --function-name $FUNCTION \
		  --region $REGION \
		  --log-type Tail \
		  --payload file://$PAYLOAD_FILE \
		  outfile_$(basename $PAYLOAD_FILE)

#jq . outfile.json

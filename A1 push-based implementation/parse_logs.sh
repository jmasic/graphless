#!/bin/bash

LOG_PARSE_START_TIME=$1

if [ -z $LOG_PARSE_START_TIME ];
then
	echo "Usage parse-logs.sh <log-parse-start-time>"
	exit 1
fi

FUNCTIONS=( MainFunction OrchestratorFunction WorkerFunction )
ID_TIME_PAIRS=$(saw get /aws/lambda/${FUNCTIONS[0]} --start $LOG_PARSE_START_TIME --filter '{ $.tag = STARTED }' | jq '.runId, .time' | tr -d \")

mkdir -p $PWD/cloudwatch_logs
cd cloudwatch_logs 

set -- $ID_TIME_PAIRS

while [ $# -gt 0 ]
do
  run_id=$1
  start_time=$2  
  end_time=$(saw get /aws/lambda/${FUNCTIONS[1]} --start $LOG_PARSE_START_TIME --filter '{ $.tag = FINISHED && $.runId ='$run_id' }' | jq '.time' | tr -d \")

  if [ -z $end_time ]; then
    echo "No end time found for $run_id with start time $start_time. Run must have failed. Skipping log parsing for this run."
    continue
  fi
  
  echo "Parsing logs for run $run_id, started at $start_time and finished at $end_time"

  run_log_dir=$PWD/$run_id
  mkdir -p $run_log_dir

  for fct_name in ${FUNCTIONS[@]}
  do
    #saw get /aws/lambda/$fct_name --start $start_time --stop $end_time >> $run_log_dir/full_logs.txt
    #saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter '{ $.runId='$run_id' }' >> $run_log_dir/monitoring_logs.json
    
    saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter '[report = REPORT, ...]' | awk \
   		'NF > 0 {print $5 " " $9}' >> $run_log_dir/$fct_name"_exec_and_billed_duration.txt"
 done

  saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter "{ $.runId='$run_id' && $.tag = 'ACTIVE_WORKERS' }" > \
    jq -s -c 'sort_by(.time)[] | .pureDuration, .superstep, .time | tr -d \"' > $run_log_dir/active_workers.json

   saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter "{ $.runId='$run_id' && $.tag = 'COLD_START' }" > \
    jq -s -c 'sort_by(.time)[] | .pureDuration, .superstep, .time | tr -d \"' > $run_log_dir/cold_starts.json

  saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter "{ $.runId='$run_id' && ($.tag = 'SEND_MESSAGE_ALL_EDGES' 
    || $.tag = 'SEND_MESSAGES' || $.tag = 'SEND_SINGLE_MESSAGE') }" > \
    jq -s -c 'sort_by(.time)[] | .pureDuration, .superstep, .time | tr -d \"' > $run_log_dir/send_messages.json


  saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter "{ $.runId='$run_id' && $.tag = 'GET_MESSAGES' }" > \
    jq -s -c 'sort_by(.time)[] | .pureDuration, .superstep, .time | tr -d \"' > $run_log_dir/get_messages.json  

  saw get /aws/lambda/$fct_name --start $start_time --stop $end_time --filter "{ $.runId='$run_id' && $.tag = 'GET_VERTICES' }" > \
    jq -s -c 'sort_by(.time)[] | .pureDuration, .superstep, .time | tr -d \"' > $run_log_dir/get_vertices.json

  # jq -s -c 'sort_by(.time)[]' $run_log_dir/monitoring_logs.json > $run_log_dir/tmp.json && mv $run_log_dir/tmp.json $run_log_dir/monitoring_logs.json
  # cat $run_log_dir/monitoring_logs.json | jq -c '. | select(.tag=="COLD_START")' > $run_log_dir/cold_starts.json
  # cat $run_log_dir/monitoring_logs.json | jq -c '. | select(.tag=="ACTIVE_WORKERS")' > $run_log_dir/active_workers.json

  echo "Done parsing logs for run $run_id"
  shift
  shift    
done

# saw get /aws/lambda/WorkerFunction --start -24h --filter START 

# saw get /aws/lambda/WorkerFunction --start -6 --filter REPORT | awk \
#     'NF > 0 {print $5 " " $9}' >> file.txt

# saw get /aws/lambda/MainFunction --start -24h --filter REPORT | awk \
# 'NF > 0 {print $5 " " $9}' >> file.txt


# saw get /aws/lambda/OrchestratorFunction --start -24h --filter REPORT | awk \
# 'NF > 0 {print $5 " " $9}' >> file.txt

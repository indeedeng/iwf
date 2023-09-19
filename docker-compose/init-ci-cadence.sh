#!/bin/bash

for run in {1..120}; do
  sleep 1
  echo "now trying to register iWF system search attributes..."
  if yes | cadence adm cl asa --search_attr_key IwfGlobalWorkflowVersion --search_attr_type 2; then
    break
  fi
done

yes | cadence adm cl asa --search_attr_key IwfExecutingStateIds --search_attr_type 1
yes | cadence adm cl asa --search_attr_key IwfWorkflowType --search_attr_type 1
yes | cadence adm cl asa --search_attr_key IwfFirstRunId --search_attr_type 1

# see https://github.com/indeedeng/iwf/blob/main/CONTRIBUTING.md#option-3-run-with-your-own-cadence-service
echo "now sleep for 60s so that all the search attributes can take effect"

sleep 70

cadence --do default domain register

tail -f /dev/null
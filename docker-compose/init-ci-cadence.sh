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
cadence --do default domain register

tail -f /dev/null
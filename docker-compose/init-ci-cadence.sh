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


for run in {1..60}; do
  sleep 1
  echo "now trying to register domain in Cadence..."
  if cadence --do default domain register | grep -q 'Domain default already registered'; then
    break
  fi
done



tail -f /dev/null
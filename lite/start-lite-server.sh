#!/bin/bash

checkExists () {
  # see https://github.com/temporalio/temporal/issues/4160
  if temporal  operator search-attribute list | grep -q "$1"; then
    return 0
  else
    return 1
fi
}

export PATH="$PATH:/root/.temporalio/bin"
temporal server start-dev --ip 0.0.0.0 --ui-ip 0.0.0.0 &
# add SAs...
echo "temporal server started..."
echo "now trying to register iWF system search attributes..."

for run in {1..60}; do
  sleep 1
  temporal  operator search-attribute  create --name IwfWorkflowType --type Keyword
  sleep 0.1
  temporal  operator search-attribute  create --name IwfGlobalWorkflowVersion --type Int 
  sleep 0.1
  temporal  operator search-attribute  create --name IwfExecutingStateIds --type KeywordList 
  sleep 0.1
  if checkExists "IwfWorkflowType" ] && checkExists "IwfGlobalWorkflowVersion" && checkExists "IwfExecutingStateIds" ] ; then
      echo "All search attributes are registered"
      break
    fi
done

/iwf/iwf-server start
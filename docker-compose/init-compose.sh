#!/bin/bash

# use new version of tctl so that we can skip the prompt
tctl config set version next

checkExists () {
  # see https://github.com/temporalio/temporal/issues/4160
  if tctl search-attribute list | grep -q "$1"; then
    return 0
  else
    return 1
fi
}

echo "now trying to register iWF system search attributes..."

for run in {1..60}; do
  sleep 1
  tctl search-attribute  create -name IwfWorkflowType -type Keyword -y
  sleep 0.1
  tctl search-attribute  create -name IwfGlobalWorkflowVersion -type Int -y
  sleep 0.1
  tctl search-attribute  create -name IwfExecutingStateIds -type Keyword -y
  sleep 0.1
  if checkExists "IwfWorkflowType" ] && checkExists "IwfGlobalWorkflowVersion" && checkExists "IwfExecutingStateIds" ] ; then
      echo "All search attributes are registered"
      break
    fi
done

tctl namespace register default

tail -f /dev/null

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
for run in {1..120}; do
  sleep 1
  tctl search-attribute  create -name IwfWorkflowType -type Keyword -y;
  sleep 0.1
  tctl search-attribute  create -name IwfWorkflowType -type Keyword -y
  sleep 0.1
  tctl search-attribute  create -name IwfGlobalWorkflowVersion -type Int -y
  sleep 0.1
  tctl search-attribute  create -name IwfExecutingStateIds -type Keyword -y
  sleep 0.1
  tctl search-attribute  create -name CustomKeywordField -type Keyword -y
  sleep 0.1
  tctl search-attribute  create -name CustomIntField -type Int -y
  sleep 0.1
  tctl search-attribute  create -name CustomBoolField -type Bool -y
  sleep 0.1
  tctl search-attribute  create -name CustomDoubleField -type Double -y
  sleep 0.1
  tctl search-attribute  create -name CustomDatetimeField -type Datetime -y
  sleep 0.1
  tctl search-attribute  create -name CustomStringField -type text -y

  if checkExists "IwfWorkflowType" ] && checkExists "IwfGlobalWorkflowVersion" && checkExists "IwfExecutingStateIds" && checkExists "CustomKeywordField" && checkExists "CustomIntField" && checkExists "CustomBoolField" && checkExists "CustomDoubleField" && checkExists "CustomDatetimeField" && checkExists "CustomStringField" ] ; then
    echo "All search attributes are registered"
    break
  fi

done


tail -f /dev/null
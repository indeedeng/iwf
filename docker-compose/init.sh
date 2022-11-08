#!/bin/bash

# use new version of tctl so that we can skip the prompt
tctl config set version next

# TODO break the loop after commands are successful
for run in {1..20}; do
  sleep 1
  tctl namespace register default || true
  tctl search-attribute  create -name IwfWorkflowType -type Keyword -y || true
  tctl search-attribute  create -name IwfGlobalWorkflowVersion -type Int -y || true
  tctl search-attribute  create -name IwfExecutingStateIds -type Keyword -y || true
done


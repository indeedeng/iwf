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


echo "After registering, it may take up 60s because of this issue. for Cadence to load the new search attributes." 
echo "If run the test too early, you may see error: \"IwfWorkflowType is not a valid search attribute key\""
echo "and the test would fail with: unknown decision DecisionType: Activity, ID: 0, possible causes are nondeterministic workflow definition code or incompatible change in the workflow definition"
sleep 65

echo "now register the domain to tell the tests that Cadence is ready"
for run in {1..60}; do
  echo "now trying to register domain in Cadence..."
  if cadence --do default domain register | grep -q 'Domain default already registered'; then
    break
  fi
  sleep 1
done

tail -f /dev/null
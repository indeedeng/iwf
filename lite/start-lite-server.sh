#!/bin/bash

checkExists () {
  # see https://github.com/temporalio/temporal/issues/4160
  if temporal  operator search-attribute list | grep -q "$1"; then
    return 0
  else
    return 1
fi
}

# Start MinIO in the background
mkdir -p /tmp/minio-data
export MINIO_ROOT_USER=minioadmin
export MINIO_ROOT_PASSWORD=minioadmin
minio server /tmp/minio-data --address ":9000" --console-address ":9001" &

# Wait for MinIO to be ready
echo "waiting for MinIO to start..."
for run in {1..30}; do
  sleep 1
  if curl -s http://localhost:9000/minio/health/ready > /dev/null 2>&1; then
    echo "MinIO is ready"
    break
  fi
done

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
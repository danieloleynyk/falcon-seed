#!/bin/bash

docker restart "$CONTAINER_NAME" > /dev/null
echo "Starting \"${CONTAINER_NAME}\" Debug Container..."

# shellcheck disable=SC1083
while [ "$(/usr/bin/docker inspect -f {{.State.Health.Status}} "${CONTAINER_NAME}")" != "healthy" ]; do
    sleep 0.1;
done;

return 0

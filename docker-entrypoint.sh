#!/bin/bash
set -e

SERVICE_NAME="$1"
shift # Remove the first argument (service name)

# Determine the default config path for the selected service
# This assumes configs are structured as /app/resources/configs/<service_name>/bootstrap.yaml
DEFAULT_CONFIG_PATH="/app/resources/configs/${SERVICE_NAME}/bootstrap.yaml"

# Check if a custom config path is provided via -conf in the arguments
# If -conf is provided, we assume the user knows what they're doing and let it pass through.
# Otherwise, we append our default config path.
HAS_CUSTOM_CONF=false
for arg in "$@"; do
  if [[ "$arg" == "-conf" ]]; then
    HAS_CUSTOM_CONF=true
    break
  fi
done

# If no custom -conf was provided, append our default
if ! $HAS_CUSTOM_CONF; then
  set -- "$@" "-conf" "$DEFAULT_CONFIG_PATH"
fi

# Execute the specified service binary
# The binary name should match the service name (e.g., /app/gateway, /app/helloworld)
exec "/app/${SERVICE_NAME}" "$@"
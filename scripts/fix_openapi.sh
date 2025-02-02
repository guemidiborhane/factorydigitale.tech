#!/bin/sh

# Define the file to be modified
file="$1"

# Use yq to update the YAML file for fields under components.schemas ending with .HttpError
yq eval -i '.components.schemas.*HttpError.properties.message = {"type": "string", "oneOf": [{"type": "string"}, {"type": "object", "additionalProperties": {"type": "string"}}]}' "$file"

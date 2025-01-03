#!/bin/bash


batch_size=10
batch=()

function requestBody() {
    printf "%s\n" "${batch[@]}" | jq -s -c '{SRNumbers: .}'
}

function lookup() {
    local data=$(requestBody)
    # https://api-portal.nyc.gov/api-details#api=nyc-311-public-api&operation=api-GetServiceRequestList-post
    curl -s -H "${OCP_AUTH_HEADER}" -H "Content-Type: application/json" -d "${data}" -X POST "https://api.nyc.gov/public/api/GetServiceRequestList" | \
        jq -c --argjson STATUSLOOKUP '{"614110001": "Open", "614110002": "In Progress", "614110000": "Cancel", "614110003": "Closed"}'  '.SRResponses[] | .Status = $STATUSLOOKUP[.Status] | .'
    sleep 1
}


# Read input line by line
while IFS= read -r line; do
  batch+=("\"${line}\"")
  
  # Check if the batch size is reached
  if [ "${#batch[@]}" -eq "$batch_size" ]; then
    lookup
    batch=() # Reset the batch
  fi
done

# Handle any remaining lines in the batch
if [ "${#batch[@]}" -gt 0 ]; then
    lookup
fi





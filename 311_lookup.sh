#!/bin/bash


function lookup() {
    local SRNUMBER=$1
    curl -s -H "${OCP_AUTH_HEADER}" -X GET "https://api.nyc.gov/public/api/GetServiceRequest?srnumber=${SRNUMBER}" | jq -c . 
}

lookup $1

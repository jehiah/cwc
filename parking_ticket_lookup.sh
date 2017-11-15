#!/bin/bash

DATA2018=https://data.cityofnewyork.us/resource/qpyv-8eyi.json
FIELDS="summons_number,registration_state,plate_id,plate_type,issue_date,violation_time,violation_location,house_number,street_name,intersecting_street,issuer_command,issuer_squad,violation_description,law_section,sub_division,violation_code"
curl --silent "${DATA2018}?plate_id=${1}&\$select=${FIELDS}"  | jq -c '.[]' | json2csv -p -k $FIELDS

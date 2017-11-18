#!/bin/bash

if [ "${1}" == "" ]; then
    echo "usage: $0 T123456C"
    echo -e "\t returns a CSV of FY 2018 tickets for license \"T123456C\""
    exit 1
fi

TICKETS_FY2018="https://data.cityofnewyork.us/resource/qpyv-8eyi.csv"
    # https://data.cityofnewyork.us/City-Government/Parking-Violations-Issued-Fiscal-Year-2018/pvqr-7yc4
TICKETS_FY2017="https://data.cityofnewyork.us/resource/ati4-9cgt.csv"
    # https://data.cityofnewyork.us/City-Government/Parking-Violations-Issued-Fiscal-Year-2017/2bnn-yakx
FIELDS="summons_number,registration_state,plate_id,plate_type,issue_date,violation_time,violation_location,house_number,street_name,intersecting_street,issuer_command,issuer_squad,violation_description,law_section,sub_division,violation_code"
FILE="tickets_${1}_$(date +%Y-%m-%d).csv"
echo "downloading tickets for ${1} into $FILE"
curl --silent "${TICKETS_FY2018}?plate_id=${1}&\$select=${FIELDS}&\$order=issue_date,violation_time" > $FILE || exit 1
echo "open $FILE"
open $FILE || exit 1

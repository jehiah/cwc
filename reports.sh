#!/bin/bash

if [ ! -f /usr/local/bin/bar_chart.py ]; then
   echo "requires data_hacks"
   echo "try: pip install data_hacks"
   exit 1
fi

echo ""
echo "Violations Cited"
ncwc --list-regulations --short | while read line; do
    reg=$(echo "$line" | awk -F, '{print $1}')
    def=$(echo "$line" | awk -F, '{print $2}')
    N=$(grep  -l "$reg" */notes.txt | wc -l | tr -d "\t" | tr -d " ")
    if [[ "$N" == "0" ]]; then
        continue
    fi
    reg=$(printf "%18s" "$reg")
    def=$(printf "%28s" "$(echo "$def" | cut -c 1-28)")
    echo -e "${N}\t${def} - $reg"
done | bar_chart.py -a -v -p

echo ""
echo "TLC Complaints by month"
find . -name 'notes.txt' | cut -c 3-8 | bar_chart.py --value-suffix=" complaints"

echo ""
echo "Days with TLC complaints by month"
find . -name 'notes.txt' | cut -c 3-10 | sort | uniq | cut -c 1-6 | bar_chart.py --value-suffix=" days"

echo ""
echo "Distribution of Complaints per day"
find . -name 'notes.txt' | cut -c 3-10 | sort | uniq -c | awk '{print $1}' | bar_chart.py --key-suffix=" complaints/day" --value-suffix=" days" -k -n

find . -type d | egrep '^\./[0-9]+_[^_]+_' | awk -F_ '{print $3}' > license_plates.log
find . -type d | egrep '^\./[0-9]+_[^_]+$' | awk -F_ '{print $2}' >> license_plates.log
echo ""
echo "# of License Plates: $(sort license_plates.log | uniq | wc -l)"
echo "# of Taxi Plates: $(egrep -c '[0-9][A-Z][0-9]{2}' license_plates.log)"
echo "# w/ Multiple Reports: $(sort license_plates.log | uniq -c | egrep -v '^\s+1\s' | wc -l) ($(sort license_plates.log | uniq -c | egrep -v '^\s+1\s' | awk '{print $2}' | tr '\n' ' '))"

echo "# tweeted: $(grep -l twitter.com */notes.txt | wc -l)"
echo ""

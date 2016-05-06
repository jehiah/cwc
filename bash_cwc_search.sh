
# put this contents in ~/.bash_profile
# usage: cwc $1
# which will search for the string $1 in any cwc report and opern that folder in Finder and Textmate
function cwc {
    local pattern=$1
    local target=$(ack -l "$pattern" ~/Documents/cyclists_with_cameras | head -1)
    if [ -z "$target" ]; then
        echo "No match for pattern \"$pattern\""
    else
        echo "found $target"
        mate $(dirname $target)
        open $(dirname $target)
    fi
}


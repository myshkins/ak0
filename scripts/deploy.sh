#!/bin/bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path

# source ./build_image.sh

# if [[ $? -ne 0 ]]; then
#   echo "error building image"
#   exit 1
# fi

# run the build image script that was created above
pushd .. >/dev/null
echo $(whoami)
./result | ./build/backend/ak0_image

scp ../build/backend/ak0_image pgum:/home/iceking/data/ak0/images/


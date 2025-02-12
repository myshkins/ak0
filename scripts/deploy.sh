##!/bin/sh
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path

# source ./build_image.sh

# if [[ $? -ne 0 ]]; then
#   echo "error building image"
#   exit 1
# fi

pushd .. >/dev/null
echo $(whoami)
./result > ./build/backend/ak0_2_image


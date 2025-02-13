#!/bin/bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path

image_name="ak0"


# build ak0 image locally
if docker ps -aq | grep "${image_name}" >/dev/null 2>&1; then
  echo "There are some existing containers. Prob wanna remove em first"
  exit 1
fi

if [[ -n $(git status --porcelain ) ]]; then
  echo "git tree not clean"
  while true; do
    read -p "want to continue? (y/n): " answer
    case ${answer:0:1} in
        y|Y ) echo "continuing..."; break;;
        n|N ) echo "exiting..."; exit 0;;
        * ) echo "please answer yes or no.";;
    esac
done
fi

if docker image ls "${image_name}" | grep "${image_name}" >/dev/null 2>&1; then
  echo "removing old docker image ${image_name}"
  docker image rm "${image_name}:latest"
fi

pushd .. >/dev/null
nix build .#docker
if [[ $? -ne 0 ]]; then
  echo "error building image"
  exit 1
fi

docker load < ./result

popd >/dev/null


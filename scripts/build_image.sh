#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


echo ""
echo "running build_image.sh"
image_name="ak0"

if docker ps -a | grep "${image_name}" >/dev/null 2>&1; then
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

./build_frontend.sh

if ! git diff --quiet; then
    echo "Changes detected after frontend build."
    echo -n "Amend the last commit? (y/n): "
    read -r response

    if [[ "$response" =~ ^[Yy]$ ]]; then
        git -C .. add .
        git commit --amend --no-edit
        git push -f
        echo "Last commit amended with build changes"
    else
        echo "Skipping commit amendment"
    fi
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

./result | docker image load

popd >/dev/null

echo "build_image.sh complete"
echo ""

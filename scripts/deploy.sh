#!/usr/bin/env bash
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
./result > ./build/backend/ak0_image

scp ./build/backend/ak0_image pgum:/home/iceking/data/ak0/images/
scp -r ./otel pgum:/home/iceking/.local/ak0/otel

ssh pgum << 'EOF'
  cd /home/iceking/apps/ak0
  docker compose -f compose.yaml -f compose.prod.yaml down
  docker image rm ak0:latest
  docker image load /home/iceking/data/ak0/images/ak0
  git pull
  ./scripts/export.sh
  docker compose up
EOF

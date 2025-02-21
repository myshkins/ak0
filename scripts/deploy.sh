#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path

echo $(pwd)
# fix error handling
./build_image.sh

if [[ $? -ne 0 ]]; then
  echo "error building image"
  exit 1
fi

# run the build image script that was created above
pushd .. >/dev/null
./result > ./build/backend/ak0_image

scp ./build/backend/ak0_image pgum:/home/iceking/data/ak0/images/
scp ./.env pgum:/home/iceking/.local/ak0/

ssh pgum << 'EOF'
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml down
  docker rm ak0_web
  docker image rm ak0:latest
  docker load -i /home/iceking/data/ak0/images/ak0_image
  cd /home/iceking/apps/ak0
  git pull
  ./configs/create_certs.sh
  ./scripts/export.sh
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF


#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

# todo: fix error handling
./build_image.sh

# run the build image script that was created above
cd .. >/dev/null
./result > ./build/backend/ak0_image

scp ./build/backend/ak0_image pgum:/home/iceking/data/ak0/images/
scp ./.env pgum:/home/iceking/.local/ak0/

ssh pgum << 'EOF'
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml down
  docker image rm ak0:latest
  docker load -i /home/iceking/data/ak0/images/ak0_image
  cd /home/iceking/apps/ak0
  git pull
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF

ssh rpgum << 'EOF'
  /home/iceking/apps/ak0/configs/create_certs.sh
  /home/iceking/apps/ak0/scripts/export.sh
EOF

ssh pgum << 'EOF'
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF

echo "deploy complete, yay"

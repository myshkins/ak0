#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

# todo: fix error handling
./build_image.sh

# run the build image script that was created above
cd .. >/dev/null
echo "running ./result to build image"
./result > ./build/backend/ak0_image

echo "copying image and .env to pgum"
scp -F /home/myshkins/.ssh/config ./build/backend/ak0_image pgum:/home/iceking/data/ak0/images/
scp -F /home/myshkins/.ssh/config ./.env pgum:/home/iceking/.local/ak0/

echo "executing deploy steps on pgum"
ssh -F /home/myshkins/.ssh/config pgum << 'EOF'
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml down
  docker image rm ak0:latest
  docker load -i /home/iceking/data/ak0/images/ak0_image
  cd /home/iceking/apps/ak0
  git pull
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF

ssh -F /home/myshkins/.ssh/config rpgum << 'EOF'
  /home/iceking/apps/ak0/configs/create_certs.sh
  /home/iceking/apps/ak0/scripts/export.sh
EOF

ssh -F /home/myshkins/.ssh/config pgum << 'EOF'
  cd /home/iceking/.local/ak0/
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF

echo "deploy complete, yay"

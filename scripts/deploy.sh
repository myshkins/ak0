#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

run_dir=/home/iceking/.local/ak0
checkout_dir=/home/iceking/apps/ak0

# todo: fix error handling
./build_image.sh

# run the build image script that was created above
cd .. >/dev/null
echo "running ./result to build image"
./result > ./build/backend/ak0_image

echo "copying image and .env to pgum"
scp -F /home/myshkins/.ssh/config ./build/backend/ak0_image pgum:/home/iceking/data/ak0/images/
scp -F /home/myshkins/.ssh/config ./.env "pgum:${run_dir}"

echo "executing deploy steps on pgum"
ssh -F /home/myshkins/.ssh/config pgum << EOF
  cd "${run_dir}"
  docker compose --profile full -f compose.yaml -f compose.prod.yaml down
  docker image rm ak0:latest
  docker load -i /home/iceking/data/ak0/images/ak0_image
  cd "${checkout_dir}"
  git pull
EOF

ssh -F /home/myshkins/.ssh/config rpgum << EOF
  "${checkout_dir}/scripts/export.sh"
  rm -rf "${run_dir}/configs/certs/*"
  "${run_dir}/configs/create_certs.sh"
EOF

ssh -F /home/myshkins/.ssh/config pgum << EOF
  cd "${run_dir}"
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF

echo "deploy complete, yay"

#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

run_dir=/home/iceking/.local/ak0
checkout_dir=/home/iceking/apps/ak0
build_image="true"
remove_volumes="false"

usage() {
  echo "Usage: ./deploy.sh [options]"
  echo "  --no-image-build    skip building the ak0 image"
  echo "  --remove-volumes    remove the docker volumes on pgum"
  exit 1
}


while [ "$1" != "" ]; do
    param=${1%%=*}
    # VALUE=${1#*=}
    case $param in
        --no-image-build)
            build_image="false"
            ;;
        --remove-volumes)
            remove_volumes="true"
            ;;
        *)
            usage
            exit 1
            ;;
    esac
    shift
done

docker_volume_flag=""
if [[ "$remove_volumes" == "true" ]];then
  echo "will remove docker volumes"
  docker_volume_flag="-v"
fi

if [[ "$build_image" == "true" ]];then
  # todo: fix error handling
  ./build_image.sh

  # run the build image script that was created above
  pushd .. >/dev/null
  echo "running ./result to build image"
  ./result > ./build/backend/ak0_image
  echo "copying image to pgum"
  scp -F /home/myshkins/.ssh/config ./build/backend/ak0_image pgum:/home/iceking/data/ak0/images/
  popd>/dev/null
fi

pushd .. >/dev/null
echo "copying .env to pgum"
scp -F /home/myshkins/.ssh/config ./.env "pgum:${run_dir}"
scp -F /home/myshkins/.ssh/config ./configs/nginx.conf rpgum:/etc/nginx/conf.d/ak0.conf

echo "executing deploy steps on pgum"

ssh -F /home/myshkins/.ssh/config pgum << EOF
  cd "${run_dir}"
  docker compose --profile full -f compose.yaml -f compose.prod.yaml down
  docker image rm ak0:latest
  docker load -i /home/iceking/data/ak0/images/ak0_image
  cd "${checkout_dir}"
  git pull
  "${checkout_dir}/scripts/export.sh"
  cd "${run_dir}"
  docker compose --profile full -f compose.yaml -f compose.prod.yaml up -d
EOF

ssh -F /home/myshkins/.ssh/config rpgum "nginx -s reload"


echo "deploy complete, yay"

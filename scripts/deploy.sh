#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path


set -e

run_dir=/home/iceking/.local/ak0
checkout_dir=/home/iceking/apps/ak0
build_image="true"
remove_volumes="false"
build_blog="false"

usage() {
  echo "Usage: ./deploy.sh [options]"
  echo "  --no-image-build    skip building the ak0 image"
  echo "  --blog-build        build the blog pages"
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
        --build-blog)
            build_blog="true"
            ;;
        *)
            usage
            exit 1
            ;;
    esac
    shift
done

if [[ $(git branch --show-current) -ne "main" ]];then
  echo "you're not on the main branch"
  exit 1
fi

# check that remote branch is up to date
git fetch origin
if ! git diff --quiet HEAD origin/main; then
    echo "local and remote branches out of sync"
    echo "do a git push"
    exit 1
    # Your code for when branches are in sync
fi

docker_volume_flag=""
if [[ "$remove_volumes" == "true" ]];then
  echo "will remove docker volumes"
  docker_volume_flag="-v"
fi

# if [[ "$build_blog" == "true" ]];then
# # todo: add this logic
# fi

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
echo "copying .env and configs to pgum"
scp -F /home/myshkins/.ssh/config ./.env "pgum:${run_dir}"
scp -F /home/myshkins/.ssh/config ./configs/nginx.conf rpgum:/etc/nginx/conf.d/ak0.conf
scp -F /home/myshkins/.ssh/config ./configs/logrotate rpgum:/etc/logrotate.d/ak0

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

ssh -F /home/myshkins/.ssh/config rpgum "chown 65532:65532 /home/iceking/.local/ak0/configs/config.json"
ssh -F /home/myshkins/.ssh/config rpgum "nginx -s reload"


echo "deploy complete, yay"

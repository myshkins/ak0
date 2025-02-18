#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path

echo ""
echo "running export.sh"
# create this hacky run location that avoids needing sudo
mkdir -p /home/iceking/.local/ak0/

pushd ./.. >/dev/null

# copy everything to run location
cp -r \
  compose.yaml \
  compose.prod.yaml \
  configs \
  grafana \
  /home/iceking/.local/ak0/

echo "export.sh complete"
echo ""

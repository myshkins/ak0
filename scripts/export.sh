#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd $script_path

# create this hacky run location that avoids needing sudo
mkdir -p /home/iceking/.local/ak0/

# copy everything to run location
cp -r \
  compose.yaml \
  compose.prod.yaml \
  configs \
  grafana \
  configs/certs \
  /home/iceking/.local/ak0/


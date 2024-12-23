#!/bin/sh
script=$(readlink -f "$0")
script_path=$(dirname "$script")

cd script_path

reflex -R 'cmd$' -R '^web/dist.*' -s -- sh -c './build_frontend.sh && cd cmd && go build && ./cmd'

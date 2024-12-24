#!/bin/sh
script=$(readlink -f "$0")
script_path=$(dirname "$script")

cd script_path

reflex -R 'cmd$' -R '.*vite.*' -R '.*dist.*' -s --verbose=true -- sh -c './build_frontend.sh && cd cmd && go build && ./cmd'

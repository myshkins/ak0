#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

cd $script_path
cd ..

reflex -R '^cmd/ak0/' -R '.*vite.*' -R '.*dist.*' -R "dev.log" -s --verbose=true -- sh -c './scripts/build_frontend.sh && cd cmd/ak0 && go build && ./ak0 --env=dev'

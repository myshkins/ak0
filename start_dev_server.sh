#!/bin/sh
script=$(readlink -f "$0")
script_path=$(dirname "$script")
echo $script_path

cd $script_path

reflex -R '^cmd/ak0_2/' -R '.*vite.*' -R '.*dist.*' -R "dev.log" -s --verbose=true -- sh -c './build_frontend.sh && cd cmd/ak0_2 && go build && ./ak0_2 --env=dev'

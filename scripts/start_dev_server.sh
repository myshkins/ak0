#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")

cd $script_path
cd ..



reflex -R '^cmd/ak0/' -R '.*vite.*' -R '.*dist.*' -R "dev.log" -R "internal/handlers/dist" -s --verbose=true -- sh -c './scripts/build_frontend.sh && cp -r web/dist/ internal/handlers/dist/ && cd cmd/ak0 && go build && ./ak0 --env=dev'

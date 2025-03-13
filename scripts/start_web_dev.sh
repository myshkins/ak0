#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd "${script_path}"

cd "../web" >/dev/null

reflex \
  -R "build" \
  -s --verbose=true \
  -- sh -c '../scripts/build_frontend.sh && cd build && npm exec live-server --no-browser'

#!/usr/bin/env bash

script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd "${script_path}"

pushd "../web" >/dev/null
npm run dev -- --config vite.config.js --open blog-dev.html
popd >/dev/null

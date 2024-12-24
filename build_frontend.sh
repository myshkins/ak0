#!/bin/sh
script=$(readlink -f "$0")
script_path=$(dirname "$script")

pushd "${script_path}/web" >/dev/null
npm run build
popd >/dev/null

#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd "${script_path}"

cd "../web" >/dev/null

rm -rf build/*
mkdir -p build/assets/
cp -r src/assets/* build/assets/
../scripts/htmlbuilder/htmlbuilder

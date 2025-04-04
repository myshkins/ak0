#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd "${script_path}"

set -e

outDir="../cmd/ak0/dist/"

go build -o htmlbuilder/htmlbuilder htmlbuilder/htmlbuilder.go
go build -o blogger/blogger blogger/blogger.go

cd "../web" >/dev/null

# go build -
# start with fresh directories
rm -rf build/*
rm -rf src/blog/posts/*
mkdir -p "${outDir}"
rm -rf "${outDir}"*
mkdir -p build/assets/
mkdir -p build/blog/

# parse md files and generate html files
../scripts/blogger/blogger
../scripts/htmlbuilder/htmlbuilder

cp -r src/assets/* build/assets/

# copy web/build to handler dir for go file embed
cp -r build/* $outDir
cp /home/myshkins/projects/job_search/resume/resume_Alex_Krenitsky.pdf "${outDir}assets/"
cp robots.txt $outDir

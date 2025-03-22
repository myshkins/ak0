#!/usr/bin/env bash
script=$(readlink -f "$0")
script_path=$(dirname "$script")
cd "${script_path}"

launch_live_server() {
  ./build_frontend.sh
  cd "../web/build"
  npm exec live-server --no-browser
}

launch_reflex() {
  cd "${script_path}/.."
  reflex \
    -R "web/build" \
    -R "cmd/ak0/dist" \
    -R "scripts" \
    -R "internal/handlers/dist" \
    -R "web/src/blog" \
    --verbose=true \
    -- sh -c './scripts/build_frontend.sh'
}

launch_blog_reflex() {
  cd /home/myshkins/coeus/coeus/blog/
  reflex -v -- sh -C "${script_path}/build_frontend.sh"
}

(trap 'kill 0' SIGINT; launch_live_server & launch_reflex & launch_blog_reflex & wait)

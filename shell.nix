{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.buildPackages; [
      go
      gopls
      nodejs_22
      reflex
    ];

    shellHook = ''
      # Install Vite
      if [ ! -d ./web/node_modules ]; then
        npm install
      fi
    '';
}

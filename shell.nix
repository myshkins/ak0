let
  pkgs = import <nixpkgs> {};
  go = pkgs.go;
  gopls = pkgs.gopls;

in pkgs.mkShell {
  packages = [
    go
    gopls
  ];

  }

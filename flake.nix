{
  description = "ak0 go backend flake";

  inputs = {
    # nixpkgs.url = "nixpkgs/nixos-21.11";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
        version = builtins.substring 0 8 lastModifiedDate;

        # Create an overlay to pin Python package versions
        pythonOverlay = final: prev: {
          python3 = prev.python3.override {
            packageOverrides = pyFinal: pyPrev: {
              markdown = pyPrev.markdown.overridePythonAttrs (old: rec {
                version = "3.6";
                src = prev.fetchPypi {
                  inherit version;
                  pname = "Markdown";
                  hash = "sha256-7U9B9trsvuuW5XbOQUxB0th22qmhbLNfqO2MLd+tAiQ=";
                };
              });
            };
          };
        };

        pkgs = import nixpkgs {
          inherit system;
          overlays = [ pythonOverlay ];
        };
        pythonEnv = pkgs.python3.withPackages (ps: with ps; [
          bcrypt
          markdown
        ]);

    in {
      packages = {
        backend = pkgs.buildGoModule {
          pname = "ak0";
          inherit version;
          src = ./.;
          vendorHash = null;
        };

        # frontend = pkgs.buildNpmPackage {
        #   pname = "ak0-frontend";
        #   version = "0.0.1";
        #   src = ./web;
        #   npmDepsHash = "sha256-InkMefNQA6e3Ul8PY8pkpXSCqaysGh10t7C683AS5LA=";
        # };

        docker = pkgs.dockerTools.streamLayeredImage {
          name = "ak0";
          tag = "latest";
          # fromImage = pkgs.dockerTools.pullImage {
          #   imageName = "alpine";
          #   imageDigest = "sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099";
          #   hash = "sha256-C3TOcLa18BKeBfS5FSe0H6BALGA/zXSwSZstK+VaPyo=";
          # };

          fromImage = pkgs.dockerTools.pullImage {
            imageName = "gcr.io/distroless/base-debian12"; #nonroot
            imageDigest = "sha256:97d15218016debb9b6700a8c1c26893d3291a469852ace8d8f7d15b2f156920f";
            hash = "sha256-p8Hmw0W3AibT1quCFMLmKO7dRMgjS8BXANPOfrQRe5g=";
          };
          contents = [
            self.packages.${system}.backend
            # self.packages.${system}.frontend
          ];
          fakeRootCommands = ''
            mkdir /ak0
            touch /ak0/ak0.log
            chmod -R 644 /ak0
            chown -R 65532:65532 /ak0"
          '';
          config = {
            User = "65532:65532";
            Cmd = [
              "/bin/ak0"
            ];
            WorkingDir = "/ak0";
            ExposedPorts = {
              "8200" = {};
            };
          };
        };
      };

      devShell = with pkgs;
        mkShell {
          buildInputs = with pkgs; [
            bash
            go
            gopls
            hugo
            nodejs_22
            nil
            pythonEnv
            reflex
          ];
          shellHook = ''
            if [ ! -d ./web/node_modules ]; then
              pushd web
              npm install
              popd
            fi
          '';
        };
      });
}


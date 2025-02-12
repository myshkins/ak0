{
  description = "ak0_2 go backend flake";

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
        pkgs = import nixpkgs { inherit system; };

    in {
      packages = {
        backend = pkgs.buildGoModule {
          pname = "ak0_2";
          inherit version;
          src = ./.;
          vendorHash = null;
        };

        frontend = pkgs.buildNpmPackage {
          pname = "ak0_2-frontend";
          version = "0.0.1";
          src = ./web;
          npmDepsHash = "sha256-InkMefNQA6e3Ul8PY8pkpXSCqaysGh10t7C683AS5LA=";
        };

        docker = pkgs.dockerTools.streamLayeredImage {
          name = "ak0_2";
          tag = "latest";
          # fromImage = pkgs.dockerTools.pullImage {
          #   imageName = "alpine";
          #   imageDigest = "sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099";
          #   hash = "sha256-C3TOcLa18BKeBfS5FSe0H6BALGA/zXSwSZstK+VaPyo=";
          # };

          fromImage = pkgs.dockerTools.pullImage {
            imageName = "gcr.io/distroless/base-debian12";
            imageDigest = "sha256:74ddbf52d93fafbdd21b399271b0b4aac1babf8fa98cab59e5692e01169a1348";
            hash = "sha256-z5xmx1oaxgxYwdEVadlRp1DmokAOounOV1gKG1o4ubI=";
          };
          contents = [
            self.packages.${system}.backend
            self.packages.${system}.frontend
            pkgs.curl
          ];
          config = {
            Cmd = [
              "/bin/ak0_2"
              "--env=prod"
            ];
            ExposedPorts = {
              "8200" = {};
            };
          };
        };
      };

      devShell = with pkgs;
        mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            nodejs_22
            nil
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


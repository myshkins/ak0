{
  description = "ak0_2 go backend flake";

  inputs = {
    # nixpkgs.url = "nixpkgs/nixos-21.11";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    # flake-utils.url = "github.numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
        version = builtins.substring 0 8 lastModifiedDate;
        pkgs = import nixpkgs { inherit system; };
        # supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
        # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
        # forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
        # Nixpkgs instantiated for supported system types.
        # nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

    in {
      packages = {
        backend = pkgs.buildGoModule {
          pname = "ak0_2";
          inherit version;
          src = ./.;
          vendorHash = null;
        };

        frontend = pkgs.stdenv.mkDerivation {
          pname = "ak0_2-frontend";
          version = "0.0.1";
          src = ./web;
          buildInputs = [ pkgs.nodejs_23 ];
          buildPhase = ''
            npm install --verbose
            npm run build
          '';
          installPhase = ''
            mkdir -p $out
            cp -r dist/* $out/
          '';
        };

        docker = pkgs.dockerTools.buildLayeredImage {
          name = "ak0_2";
          tag = "latest";
          fromImage = pkgs.dockerTools.pullImage {
            imageName = "gcr.io/distroless/base-debian12";
            imageDigest = "sha256:74ddbf52d93fafbdd21b399271b0b4aac1babf8fa98cab59e5692e01169a1348";
            hash = "sha256-z5xmx1oaxgxYwdEVadlRp1DmokAOounOV1gKG1o4ubI=";
          };
          contents = [
            self.packages.${system}.backend
            self.packages.${system}.frontend
          ];
          config = {
            Cmd = [ "/ak0_2"];
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
              npm install
            fi
          '';
        };
      });
}


{
  description = "ak0_2 go backend flake";

  inputs = {
    # nixpkgs.url = "nixpkgs/nixos-21.11";
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    utils.url = "github:numtide/flake-utils";
    # flake-utils.url = "github.numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }: 
    utils.lib.eachDefaultSystem (system:
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
        buildInputs = [ pkgs.nodejs ];
        buildPhases = ''
          npm install
          npm run build
        '';
        installPhase = ''
          mkdir -p $out
          cp -r dist/* $out/
        '';
      };
      
      docker = {
        name = "ak0_2";
        tag = "latest";
        contents = [
          self.packages.${system}.backend
          self.packages.${system}.frontend
        ];
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


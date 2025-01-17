{
  description = "ak0_2 go backend flake";

  inputs.nixpkgs.url = "nixpkgs/nixos-21.11";
  # inputs.nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";

  outputs = { self, nixpkgs }: 
    let

      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
      version = builtins.substring 0 8 lastModifiedDate;

      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      # pkgs = import nixpkgs { system = "x86_64-linux"; };

    in {
      
      packages = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          go-backend = pkgs.buildGoModule {
            pname = "ak0_2";
            inherit version;
            src = ./.;
            vendorHash = null;
          };  
        }
      );

      devShells = forAllSystems (system:
        let pkgs = nixpkgsFor.${system};
        in {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools ];
          };
        });
    };
}


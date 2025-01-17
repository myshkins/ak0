{
  description = "ak0_2 go backend flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: {
    packages.x86_64-linux.default =
    let
      pkgs = import nixpkgs { system = "x86_64-linux"; };
    in pkgs.buildGoModule {
      pname = "ak0_2";
      version = "0.0.1";
      src = ./.;
      venderSha256 = null;
    };
  };
}

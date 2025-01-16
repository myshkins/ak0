{
  # lib,
  # stdenv,
  fetchFromGithub,
  buildGoModule,
}:

let
  version = "0.0.1";
in
buildGoModule {
  pname = "ak0_2-go-backend";
  inherit version;

  src = fetchFromGithub {
    owner = "myshkins";
    repo = "ak0_2";
    rev = "v${version}";
    sha256 = "sha256-77486e9b7c7ba7f54b0ea11bd6f0319350380a02=";
  };

  vendorSha256 = "";
}

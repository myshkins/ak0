{
  pkgs,
  buildLayeredImage,
  pullImage,
}:

let
  version = "0.0.1";
  appFlake = import ./flake.nix;
in

rec {
  baseImage = pullImage {
    imageName = "";
    imageDigest = "";
    hash = "";
    finalImageTag = "";
    finalImageName = "";
  };

  ak0_2 = buildLayeredImage {
    name = "ak0_2";
    tag = "v${version}";
    fromImage = baseImage;
    contents = [ appFlake ];
  };
}
 

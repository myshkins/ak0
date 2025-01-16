{
  pkgs,
  buildImage,
  buildLayeredImage,
  pullImage,
}:

let
  frontendArgs =
    let
      version = "0.0.1";
      nginxConf = pkgs.writeText "nginx.conf" "";
    in
    {
      name = "frontendImage";
      tag = "latest";
      contents = [
        pkgs.nginx
      ];

      extraCommands = ''
        mkdir -p tmp/nginx_client_body
        mkdir -p var/log/nginx
      '';

      config = {
        Cmd = [
          "nginx"
          "-c"
          nginxConf
        ];
        ExposedPorts = {
          "${nginxPort}/tcp" = { };
        };
      };
    };

  nginxArguments =
    let
      nginxPort = "80";
      nginxConf = pkgs.writeText "nginx.conf" "";
    in
    {
      name = "nginx-container";
      tag = "latest";
      contents = [
        pkgs.nginx
      ];

      extraCommands = ''
        mkdir -p tmp/nginx_client_body
        mkdir -p var/log/nginx
      '';

      config = {
        Cmd = [
          "nginx"
          "-c"
          nginxConf
        ];
        ExposedPorts = {
          "${nginxPort}/tcp" = { };
        };
      };
    };

in

rec {
  # 1. basic example
  bash = buildImage {
    name = "bash";
    tag = "latest";
    copyToRoot = pkgs.buildEnv {
      name = "image-root";
      paths = [ pkgs.bashInteractive ];
      pathsToLink = [ "/bin" ];
    };
  };

  # 2. service example, layered on another image
  redis = buildImage {
    name = "redis";
    tag = "latest";

    # for example's sake, we can layer redis on top of bash or debian
    fromImage = bash;
    # fromImage = debian;

    copyToRoot = pkgs.buildEnv {
      name = "image-root";
      paths = [ pkgs.redis ];
      pathsToLink = [ "/bin" ];
    };

    runAsRoot = ''
      mkdir -p /data
    '';

    config = {
      Cmd = [ "/bin/redis-server" ];
      WorkingDir = "/data";
      Volumes = {
        "/data" = { };
      };
    };
  };

  # 3. another service example
  backend = buildLayeredImage {};
  frontend = buildLayeredImage {};

  # Used to demonstrate how virtualisation.oci-containers.imageStream works
  nginxStream = pkgs.dockerTools.streamLayeredImage nginxArguments;

  # 4. example of pulling an image. could be used as a base for other images
  nixFromDockerHub = pullImage {
    imageName = "nixos/nix";
    imageDigest = "sha256:85299d86263a3059cf19f419f9d286cc9f06d3c13146a8ebbb21b3437f598357";
    hash = "sha256-xxZ4UW6jRIVAzlVYA62awcopzcYNViDyh6q1yocF3KU=";
    finalImageTag = "2.2.1";
    finalImageName = "nix";
  };


  # 10. Create a layered image
  ak0_2Image = pkgs.dockerTools.buildLayeredImage {
    name = "ak0_2Image";
    tag = "latest";
    extraCommands = ''echo "(extraCommand)" > extraCommands'';
    config.Cmd = [ "${pkgs.hello}/bin/hello" ];
    contents = [
      backend
      frontend
      pkgs.go
      pkgs.nodejs_22
    ];
  };
}

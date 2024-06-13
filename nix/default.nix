{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ../flake.lock)).nodes) nixpkgs gomod2nix;
    in
      import (fetchTree nixpkgs.locked) {
        overlays = [
          (import "${fetchTree gomod2nix.locked}/overlay.nix")
        ];
      }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
  version,
}:
buildGoApplication {
  pname = "site";
  inherit version;
  pwd = ../.;
  src = ../.;
  subPackages = ["cmd/site"];
  modules = ../gomod2nix.toml;
  ldflags = [
    "-w"
    "-s"
    "-X main.version=${version}"
  ];
  postInstall = ''
    mkdir -p $out/bin/env
    cp -R static $out/bin
    cp env/sample.config $out/bin/env/config
    rm $out/bin/static/resume.tex
  '';
}

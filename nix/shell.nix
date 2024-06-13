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
  mkGoEnv ? pkgs.pkgs.mkGoEnv,
  gomod2nix ? pkgs.gomod2nix,
}: let
  goEnv = mkGoEnv {
    pwd = ../.;
  };
in
  pkgs.mkShell {
    packages = [
      pkgs.go # go_1_22
      pkgs.nodejs_22
      pkgs.gotools
      pkgs.gofumpt
      pkgs.templ
      pkgs.air

      goEnv
      gomod2nix

      ## asset gen
      pkgs.texliveFull
    ];
  }

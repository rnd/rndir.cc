{
  description = "A flake for rndir.cc website";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    gomod2nix,
  }: (
    flake-utils.lib.eachDefaultSystem
    (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.default = pkgs.callPackage ./nix/default.nix {
        inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
        version = builtins.substring 0 8 self.lastModifiedDate;
      };
      devShells.default = pkgs.callPackage ./nix/shell.nix {
        inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
      };

      ##nixosModules.default = {
      ##  self,
      ##  pkgs,
      ##  ...
      ##}: {
      ##  systemd.services.site = {
      ##    description = "github.com/rnd/site";
      ##    wantedBy = ["multi-user.target"];
      ##    serviceConfig = {
      ##      Type = "simple";
      ##      ExecStart = "${self.package.${pkgs.system}.default}/bin/site";
      ##    };
      ##  };
      ##};
    })
  );
}

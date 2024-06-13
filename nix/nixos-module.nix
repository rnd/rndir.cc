self: {
  config,
  lib,
  ...
}:
with lib; let
  cfg = config.services.rnd-site;
in {
  options.services.rnd-site = {
    enable = mkEnableOption "Activate rndir.cc website";

    port = mkOption {
      type = types.port;
      default = 7001;
      example = 7001;
      description = "The port number should listen on";
    };
  };

  config = mkIf cfg.enable {
    systemd.services.rnd-site = {
      wantedBy = ["multi-user.target"];

      serviceConfig = {
        Type = "simple";
        ExecStart = "${self.package.${pkgs.system}.default}/bin/site";
      };
    };

    services.nginx.virtualHosts."rndir.cc" = {
      locations."/" = {
        proxyPass = "http://127.0.0.1:${toString cfg.port}";
      };
    };
  };
}

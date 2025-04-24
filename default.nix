{ pkgs, name, version, ... }:
pkgs.buildGoApplication {
  pname = name;
  version = version;

  src = builtins.path {
    path = ./.;
    name = "source";
  };
}

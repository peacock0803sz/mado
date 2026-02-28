{
  description = "";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];

      flake = {
        homeManagerModules.default = ./nix/hm-module.nix;
        darwinModules.default = ./nix/darwin-module.nix;
      };
      perSystem = { pkgs, lib, ... }: {
        packages = lib.optionalAttrs pkgs.stdenv.isDarwin {
          default = pkgs.callPackage ./nix/package.nix { };
        };

        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            git
            go
            golangci-lint
            gofumpt
          ];
        };
      };
    };
}

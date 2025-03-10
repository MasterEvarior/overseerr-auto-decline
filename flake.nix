{
  description = "Development flake for overseer-auto-decline";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
  };

  outputs =
    { nixpkgs, ... }:
    let
      x86 = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages."${x86}";
    in
    {
      devShells."${x86}".default = pkgs.mkShell {
        packages = with pkgs; [
          go

          # Formatters
          treefmt
          beautysh
          golangci-lint
          mdformat
          yamlfmt
          deadnix
          nixfmt-rfc-style
        ];

        shellHook = ''
          echo "Started overseerr-auto-decline dev shell"
        '';

        # Environment Variables
        URL = "https://your-overseer-instance";
        API_KEY = "ZHVtbXkta2V5Cg=="; # This is just a dummy for obvious reasons
        MEDIA = "8966,24021";
        #DELETE_REQUESTS = "true";
      };
    };
}

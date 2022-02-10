{
  description = "Frittenbude Exploration Stuff";
  inputs.nixpkgs.url = github:NixOS/nixpkgs;

  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
    lib = pkgs.lib;
  in {
    devShell.${system} = pkgs.mkShell rec {
      buildInputs = with pkgs; [
        go
        gopls
        stdenv
        glibc.static
      ];

      CFLAGS="-I${pkgs.glibc.dev}/include";
      LDFLAGS="-L${pkgs.glibc}/lib";
    };
  };
}

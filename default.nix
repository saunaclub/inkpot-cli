{ lib, buildGoModule, fetchFromGitHub, stdenv, glibc }:

buildGoModule rec {
  pname = "inkpot-cli";
  version = "0.0.1";

  src = ./.;
  vendorSha256 = "sha256-W+oAjjRYXoKM20nubO0y2yUA4WRjOn7zki3pIf9TMvc=";

  buildInputs = [
    stdenv
    glibc.static
  ];
  ldflags = "-linkmode external -extldflags -static";

  meta = with lib; {
    description = "Command-line tool to customize Spotify client";
    homepage = "https://github.com/khanhas/spicetify-cli/";
    license = licenses.gpl3Plus;
    maintainers = with maintainers; [ jonringer ];
  };
}

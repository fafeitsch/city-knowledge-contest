{ lib, stdenv, buildGoModule }:

buildGoModule rec {
  name = "gitea-pull-request-create-plugin";

  src = ./.;
  vendorSha256 = null;

  CGO_ENABLED = 0;

  meta = {
    description = "Woodpecker CI plugin for creating Pull Requests in Gitea";
    homepage = "https://codeberg.org/JohnWalkerx/gitea-pull-request-create-plugin";
    license = licenses.mit;
    maintainers = [ "johnwalkerx@mailbox.org" ];
    platforms = lib.platforms.linux;
  };
}
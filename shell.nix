with import <nixpkgs> {};
pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.go_1_14
    pkgs.postgresql
  ];
}

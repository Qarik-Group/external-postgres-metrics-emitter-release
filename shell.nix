with import <nixpkgs> {};
pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.go
  ];
}

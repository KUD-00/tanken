{ pkgs ? import <nixpkgs> {} }:

let
  lib = pkgs.lib;
in
pkgs.mkShell (with pkgs; {
  buildInputs = [
    git
    gh

    nodejs
    nodePackages_latest.pnpm

    go
    go-migrate
    
    protobuf
    protoc-gen-go
    protoc-gen-go-grpc
    grpcurl
  ];
})

#!/usr/bin/env bash
set -e
set -o pipefail
shopt -s nullglob

Package=github.com/khlieng/dispatch
BuildDir=$GOPATH/src/$Package/build
ReleaseDir=$GOPATH/src/$Package/release
BinaryName=dispatch

mkdir -p $BuildDir
cd $BuildDir
rm -f dispatch*
gox -ldflags -w $Package

mkdir -p $ReleaseDir
cd $ReleaseDir
rm -f dispatch*
for f in $BuildDir/*
do
  zipname=$(basename ${f%".exe"})
  if [[ $f == *"linux"* ]] || [[ $f == *"bsd"* ]]; then
    zipname=${zipname}.tar.gz
  else
    zipname=${zipname}.zip
  fi

  binbase=$BinaryName
  if [[ $f == *.exe ]]; then
    binbase=$binbase.exe
  fi
  bin=$BuildDir/$binbase
  mv $f $bin

  if [[ $zipname == *.zip ]]; then
    zip -j $zipname $bin
  else
    tar -cvzf $zipname -C $BuildDir $binbase
  fi

  mv $bin $f
done

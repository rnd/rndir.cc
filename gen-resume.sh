#!/usr/bin/env bash

docker run --rm -i -v "$PWD":/data latex xelatex \
  -output-directory=.build \
  src/resume.tex 

cp .build/resume.pdf static

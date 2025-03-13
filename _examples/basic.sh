#!/bin/bash
rm -rf basic 
mkdir -p basic
cd basic
go mod init basic

# no docs, no extras, use viper, set a repo and an author
starter init -a "Mister Smith" -l MIT --viper -r https://github.com/bketelsen/basic
go mod tidy

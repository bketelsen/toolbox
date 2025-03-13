#!/bin/bash
rm -rf everything 
mkdir -p everything
cd everything
go mod init everything

# all docs, all extras, use viper, set a repo and an author
starter init -d -e -a "Brian Smith" -l MIT --viper -r https://github.com/bketelsen/everything
go mod tidy

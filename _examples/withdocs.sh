#!/bin/bash
rm -rf withdocs 
mkdir -p withdocs
cd withdocs
go mod init withdocs

# no docs, no extras, use viper, set a repo and an author
starter init -a "Brian Smith" -l MIT --viper -r https://github.com/bketelsen/withdocs
go mod tidy
starter docs -b "/withdocs"
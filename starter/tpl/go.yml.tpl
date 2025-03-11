# This workflow will build a go project
# For more information see: https://docs.github.com/en/actions/automating-builds-and-tests/building-and-testing-go

name: Go

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:

  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    # uncomment if you need to do extra installations
    #- name: Tools
    #  run: go install github.com/a-h/templ/cmd/templ@latest
    - name: Build
      run: go build -v ./...

    - name: Test
      run: go test -v ./...

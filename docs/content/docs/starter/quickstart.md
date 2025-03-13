---
title: 'Quick Start'
weight: 2
---

Install starter, then open a terminal.

Create a new directory for your application:

```bash
mkdir basic
```

Change into that directory:

```bash
cd basic
```

Initialize a Go module for your project.

```bash
go mod init basic
```

Now use `starter` to generate an application skeleton for you:

```bash
starter init -a "Mister Smith" -l MIT --viper -r https://github.com/bketelsen/basic
```

Finally, run `go mod tidy` to install all dependencies:

```bash
go mod tidy
```

You've not got a starter application that doesn't (yet) do anything.

You can run it:

```
go run main.go
```

You should see the help output of the root command and some other information:

```
basic version unknown (unknown)
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  basic [flags]

Flags:
      --config string   config file (default is $HOME/.basic.yaml)
  -h, --help            help for basic
  -v, --verbose         verbose logging
      --version         version for basic
```

### Available Tasks

#### build - Build the application

Builds the application and embeds the version information in the binary



Run this task:

```
task build
```

#### checks - Run all go checks

Runs staticcheck, vet, test, format, and tidy


Run this task:

```
task checks
```

#### direnv - Add direnv hook to your bashrc


Run this task:

```
task direnv
```

#### format - Format all Go source

Format all Go source


Run this task:

```
task format
```

#### goreleaser - Install goreleaser on debian derivatives


Run this task:

```
task goreleaser
```

#### publish - Push and tag at .VERSION

Push and tag at taskfile variable "VERSION", which triggers the release process on GitHub. 
Update the version in the taskfile and commit the change before running this task.



Run this task:

```
task publish
```

#### release-check - Run goreleaser check

Run goreleaser check to verify the release configuration



Run this task:

```
task release-check
```

#### snapshot - Run goreleaser in snapshot mode

Run goreleaser in snapshot mode, suitable for testing the release process



Run this task:

```
task snapshot
```

#### staticcheck - Run go staticcheck

Run go staticcheck


Run this task:

```
task staticcheck
```

#### test - Run all tests

Run all tests


Run this task:

```
task test
```

#### tidy - Run go mod tidy

Run go mod tidy


Run this task:

```
task tidy
```

#### tools - Install required tools

Installs the `starter` command line tool


Run this task:

```
task tools
```

#### vet - Run go vet on sources

Run go vet on sources


Run this task:

```
task vet
```


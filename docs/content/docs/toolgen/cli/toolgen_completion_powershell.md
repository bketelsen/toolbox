---
date: 2025-03-21T18:00:51Z
title: "toolgen completion powershell"
slug: toolgen_completion_powershell
url: /docs/toolgen/cli/toolgen_completion_powershell/
---
## toolgen completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	toolgen completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
toolgen completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -a, --author string    author name for copyright attribution (default "YOUR NAME")
      --config string    use config file (default "toolgen.yaml")
  -l, --license string   name of license for the project
  -v, --verbose          verbose logging
```

### SEE ALSO

* [toolgen completion](/toolbox/docs/toolgen/cli/toolgen_completion/)	 - Generate the autocompletion script for the specified shell

###### Auto generated by toolbox on 21-Mar-2025

---
date: 2025-03-13T18:31:58Z
title: "starter completion powershell"
slug: starter_completion_powershell
url: /docs/starter/cli/starter_completion_powershell/
---
## starter completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	starter completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
starter completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -a, --author string       author name for copyright attribution (default "YOUR NAME")
      --config string       config file (default is starter.yaml)
  -l, --license string      name of license for the project
  -r, --repository string   gitHub repository URL (default "https://github.com/you/project")
      --verbose             verbose output
      --viper               use Viper for configuration
```

### SEE ALSO

* [starter completion](/toolbox/docs/starter/cli/starter_completion/)	 - Generate the autocompletion script for the specified shell

###### Auto generated by toolbox on 13-Mar-2025

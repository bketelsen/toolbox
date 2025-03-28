---
date: 2025-03-21T17:47:34Z
title: "thewhiz completion bash"
slug: thewhiz_completion_bash
url: /docs/cli/thewhiz_completion_bash/
---
## thewhiz completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(thewhiz completion bash)

To load completions for every new session, execute once:

#### Linux:

	thewhiz completion bash > /etc/bash_completion.d/thewhiz

#### macOS:

	thewhiz completion bash > $(brew --prefix)/etc/bash_completion.d/thewhiz

You will need to start a new shell for this setup to take effect.


```
thewhiz completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [thewhiz completion](/thewhiz/docs/cli/thewhiz_completion/)	 - Generate the autocompletion script for the specified shell

###### Auto generated by toolbox on 21-Mar-2025

<h1 align="center">flawa</h1>

Generate docs for your code with ease

> [!IMPORTANT]  
> If you are struggling with setting up something you can open an issue or DM me directly [**@0x_syk**](https://x.com/0x_syk) on x dot com the everything app.

---

- Install Golang
- Install Ollama
- `ollama pull llama3.2` and `ollama serve`

```sh
git clone https://github.com/0x-kys/flawa.git
# OR
git clone git@github.com:0x-kys/flawa.git

cd flawa && chmod +x setup.sh && ./setup.sh
```

### Usage:

```sh
A command-line tool to generate docs from your code files

Usage:
  flawa [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate document of specific file
  help        Help about any command
  list        Get info of a specific repo by name

Flags:
  -h, --help   help for flawa

Use "flawa [command] --help" for more information about a command.
```

```
flawa list <path/to/directory>

flawa generate <path/to/directory>
```

---

TODO:

- [ ] Save response to a document same as filename
- [ ] Pass entire directory and read files one by one


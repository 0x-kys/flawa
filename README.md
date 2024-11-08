<h1 align="center">flawa</h1>

Generate docs for your code with ease

---

- Install Golang
- Install Ollama
- `ollama pull llama3.2` and `ollama serve`

```sh
git clone https://github.com/0x-kys/flawa.git
# OR
git clone git@github.com:0x-kys/flawa.git

cd flawa && chmod +x setup.sh && ./setup.sh install

# to uninstall
./setup.sh uninstall
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

> Reminder to self: check todos on tiddly


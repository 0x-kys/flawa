# Flawa

**Flawa** is a command-line tool designed to generate documentation from your code files, making it easier to understand and maintain projects. It comes with commands for listing files, generating documentation, and supporting autocompletion in various shells.

## Installation

To get started, ensure you have **Golang** and **Ollama** installed, then set up Flawa with the following steps:

1. **Install Golang**: [Golang installation guide](https://golang.org/doc/install)
2. **Install Ollama**: Follow instructions at [Ollama](https://ollama.com/download)
3. **Pull and Serve Ollama**:
   ```sh
   ollama pull qwen2.5-coder # or llama3.2
   ollama serve
   ```

4. **Clone the Repository**:
   ```sh
   git clone https://github.com/0x-kys/flawa.git
   # OR
   git clone git@github.com:0x-kys/flawa.git
   ```

5. **Run Setup**:
   ```sh
   qwen2.5-coderh
   cd flawa
   chmod +x setup.sh
   ./setup.sh install
   ```

   > **To Uninstall**: Run `./setup.sh uninstall`

## config.toml

```toml
[ignore]
directories = [".git", "node_modules", "venv"] # Directories to ignore
files = ["anyrandomfile"] # Files to ignore

[ollama]
model = "qwen2.5-coder" # Choose your preferred model ()
stream = false # Keep this as false
base_prompt = '''WORKS BEST WITH DEFAULT''' # Default base prompt for the model
```

## Usage

> [!WARNING]
> **Disclaimer:** This tool utilizes pre-trained large language models (LLMs) with a general-purpose base prompt, which may result in variable or imperfect output quality. It’s recommended to thoroughly review generated documents to ensure accuracy and relevance before using them in any formal or critical context.

### General Usage

Once installed, you can access Flawa with the `flawa` command. Here’s an overview:

```sh
flawa
```

This shows a brief description, usage, available commands, and flags:

```
A command-line tool to generate docs from your code files

Usage:
  flawa [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate documentation for a specific file
  help        Help about any command
  list        List files in the specified directory or current working directory if no directory is provided

Flags:
  -h, --help   help for flawa

Use "flawa [command] --help" for more information about a command.
```

### Commands

#### `list`

The `list` command displays the files in the specified directory. If no directory is provided, it lists files in the current working directory.

```sh
# List files in the current directory
flawa list

# List files in a specific directory
flawa list path/to/directory
```

#### `generate`

The `generate` command creates documentation from code files. You can generate documentation for a single file or (in development) all files in a directory.

```sh
# Generate documentation for a specific file
flawa generate -f path/to/file

# Generate documentation for all files in a specific directory (under development)
flawa generate -d path/to/directory
```

## Contributing

Feel free to submit issues or pull requests if you have ideas for improvements or spot any bugs!

## To-Do

- [ ] [TEST] Nested subdirectories for `-d`
- [ ] [REDACTOR] `cmd/generate.go` to use single function with conditionals
- [ ] [FEATURE] Add markdown reading support (charmbracelet/glamour)
    - `flawa read path/to/flawafied-file.md`

> **Note to Self**: Check todos on Tiddly.


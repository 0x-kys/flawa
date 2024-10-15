<h1 align="center">flawa</h1>

A cli tool to interact with Github from terminal. This is something similar to gh-cli but customized to fit my workflow and improve my terminal experience.

### Usage

```txt
A command-line tool to generate docs from your code files

Usage:
  flawa [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gen         Generate docs
  help        Help about any command
  info        Get info of a specific repo by name
  login       Log in to your GitHub account using Device Flow
  repos       List all your GitHub repositories (public and private)

Flags:
  -h, --help   help for flawa

Use "flawa [command] --help" for more information about a command.
```

### Features

- Login: `flawa login`

- List Repos: `flawa repos`

  ```txt
  Repo Name                                | Visibility   | Created At   | Updated At
  ----------------------------------------------------------------------------------------------------
  public-repo                              | public       | 2024-08-14   | 2024-10-12
  private-repo                             | private      | 2024-09-06   | 2024-10-13
  ```

- Show info of a repo: `flawa info <repo-name>`

  ```txt
  Repo Name: private-repo
  Total Commits: 4
  Created At: 2024-10-11
  Updated At: 2024-10-14
  Additions: 780
  Deletions: 12
  Visibility: private
  ```

- Generate docs a specific file: `flawa gen <repo-name> --file <path-to-file>`
  ```txt
  File found, reading chunks.
  doc generated successfully
  ```

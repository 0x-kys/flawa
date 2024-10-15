<h1 align="center">flawa</h1>

Generate docs for your code with ease

> [!IMPORTANT]  
> This tool is very very alpha. If you are struggling with setting up something you can open an issue or DM me directly [**@0x_syk**](https://x.com/0x_syk) on x dot com the everything app.

> [!CAUTION]  
> This tool, for now, works with `git init`'ed repos only
>
> If your directory doesn't contain a `.git` folder it'll throw an error.
> In that case it'll try to clone repo (by using folder name) and username (from your login)
>
> If you use git ssh for cloning then you can use `--ssh` flag it'll switch cloning command to use `git@github.com`
>
> If you use an alias for github host then you can use `--domain git@yourdomain`

---

- [Setup](#setup)
- [Usage](#usage)
- [Features](#features)
- [Flags](#flags)

---

### Setup

- Golang ofc
- `wget -O - https://raw.githubusercontent.com/0x-kys/flawa/refs/heads/main/setup.sh | bash`
- and follow [these instructions](./docs/setup.md) after running the wget command

```sh
# useful when messing around!
#
# if you edit something in code
# you can build easily :3

chmod +x build.sh && ./build.sh
```

> [!NOTE]  
> `setup.sh` is for ppl who want to enjoy the `wget` way of installation
>
> `build.sh` is for ppl who want already cloned repo and just want to setup easily

---

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

---

### Flags

TODO~

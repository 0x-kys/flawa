<h1 align="center">flawa</h1>

Generate docs for your code with ease

> [!IMPORTANT]  
> This tool is very very alpha. If you are struggling with setting up something you can open an issue or DM me directly [**@0x_syk**](https://x.com/0x_syk) on x dot com the everything app.

---

- Install Golang
- Install Ollama
- `ollama pull llama3.2` and `ollama serve`

```

```

---

> [!IMPORTANT]
> create a `~/.config/flawa/config.toml` and add 

```toml
[ignore]
directories = [".git", "node_modules", "venv"]
files = ["anyrandomfile"]

[ollama]
model = 'llama3.2'
stream = false
```

```sh
flawa list <path-to-directory>

flawa generate <path-to-file>
```


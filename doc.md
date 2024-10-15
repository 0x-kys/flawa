# Imports
## External libraries
* `sh`: The shell command processor.

# go build
- **Build the Go application**:
    ```bash
go build . || handle_error "Failed to build the Go application."
```
- **Set executable permissions**:
    ```bash
chmod +x "$binary_name" || handle_error "Failed to set executable permissions."
```
- **Copy binary file**:
    ```bash
sudo cp -r "$binary_name" /usr/bin/ || handle_error "Failed to copy the binary to /usr/bin."
```

# directory creation and file checks
## Check if config directory exists
```bash
if [ ! -d "$cdir" ]; then
    mkdir -p "$cdir" || handle_error "Failed to create config directory."
else
    echo "Config directory already exists. Skipping creation."
fi
```
## Check for .env file existence
```bash
if [ ! -f "$cdir/.env" ]; then
    cp -r .env.example "$cdir/.env" || handle_error "Failed to create .env file."
else
    echo ".env file already exists. Skipping copy."
fi
```
## Check for .token file existence
```bash
if [ ! -f "$cdir/.token" ]; then
    touch "$cdir/.token" || handle_error "Failed to create .token file."
else
    echo ".token file already exists. Skipping creation."
fi
```

# Setup completion
## Output success message
```bash
echo "Setup completed successfully!"
```
## OAuth application creation reminder
```bash
echo "Finally, create an OAuth application on github.com and get your client ID and client secret."
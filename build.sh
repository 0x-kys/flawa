#!/bin/sh

cdir="$HOME/.config/flawa"
binary_name="flawa"

handle_error() {
    echo "Error: $1"
    echo "You are on your own here..."
    echo ""
    echo "Raising issue is always an option tho..."
    echo "You can ask 0x_syk on x dot com the everything app ;3"
    exit 1
}

go build . || handle_error "Failed to build the Go application."

chmod +x "$binary_name" || handle_error "Failed to set executable permissions."

if [ ! -f "/usr/bin/$binary_name" ]; then
    sudo cp -r "$binary_name" /usr/bin/ || handle_error "Failed to copy the binary to /usr/bin."
else
    echo "Binary already exists in /usr/bin. Skipping copy."
fi

if [ ! -d "$cdir" ]; then
    mkdir -p "$cdir" || handle_error "Failed to create config directory."
else
    echo "Config directory already exists. Skipping creation."
fi

if [ ! -f "$cdir/.env" ]; then
    cp -r .env.example "$cdir/.env" || handle_error "Failed to create .env file."
else
    echo ".env file already exists. Skipping copy."
fi

if [ ! -f "$cdir/.token" ]; then
    touch "$cdir/.token" || handle_error "Failed to create .token file."
else
    echo ".token file already exists. Skipping creation."
fi

echo "Setup completed successfully!"
echo "Finally, create an OAuth application on github.com and get your client ID and client secret."
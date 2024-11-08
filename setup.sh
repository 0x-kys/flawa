#!/bin/sh

set -e

if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    exit 1
fi

SCRIPT_DIR=$(dirname "$0")
cd "$SCRIPT_DIR"

install() {
    echo "Building the CLI application..."
    go build -o flawa .

    if [ ! -f "flawa" ]; then
        echo "Build failed. Could not create the binary."
        exit 1
    fi

    echo "Installing the binary to /usr/local/bin..."
    sudo cp flawa /usr/local/bin/

    if ! command -v flawa &> /dev/null; then
        echo "Installation failed. Could not find the binary in your PATH."
        exit 1
    fi

    echo "Installation successful! Use 'flawa' from anywhere."
    echo "Make sure you have /usr/local/bin in your path."

    CONFIG_DIR="$HOME/.config/flawa"
    mkdir -p "$CONFIG_DIR" 
    echo "Copying config.toml to $CONFIG_DIR"
    cp config.toml "$CONFIG_DIR/"

    if [ ! -f "$CONFIG_DIR/config.toml" ]; then
        echo "Failed to copy config.toml to $CONFIG_DIR"
        exit 1
    fi

    echo "config.toml successfully copied to $CONFIG_DIR"
}

uninstall() {
    echo "Uninstalling the binary from /usr/local/bin..."
    sudo rm -f /usr/local/bin/flawa

    if command -v flawa &> /dev/null; then
        echo "Uninstallation failed. Binary still exists."
        exit 1
    fi

    echo "Uninstallation successful."

    CONFIG_DIR="$HOME/.config/flawa"
    if [ -d "$CONFIG_DIR" ]; then
        echo "Removing $CONFIG_DIR"
        rm -rf "$CONFIG_DIR"
    fi
}

case "$1" in
    install)
        install
        ;;
    uninstall)
        uninstall
        ;;
    *)
        echo "Usage: $0 {install|uninstall}"
        exit 1
        ;;
esac


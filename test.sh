#!/bin/sh

# Set variables
cdir="$HOME/.config/flawa"
binary_name="flawa"

# Function to handle errors
handle_error() {
    echo "Error: $1"
    echo "You are on your own here..."
    echo ""
    echo "Raising issue is always an option tho..."
    echo "You can ask 0x_syk on x dot com the everything app ;3"
    exit 1
}

# Build the Go application
go build . || handle_error "Failed to build the Go application."

# Make the binary executable
chmod +x "$binary_name" || handle_error "Failed to set executable permissions."

# Copy the binary to /usr/bin
sudo cp -r "$binary_name" /usr/bin/ || handle_error "Failed to copy the binary to /usr/bin."

# Create the config directory if it doesn't exist
mkdir -p "$cdir" || handle_error "Failed to create config directory."

# .env and .token files
cp -r .env.example "$cdir/.env" || handle_error "Failed to create .env file."
touch "$cdir/.token" || handle_error "Failed to create .token file."

echo "Setup completed successfully!"

echo "Finally create OAuth application on github dot com and get your client ID and client secret."
#!/bin/bash

# Define the name of the binary and the URL to download it from
BINARY_NAME="pinata-go-cli"
BINARY_URL="https://github.com/stevedylandev/pinata-go-cli/pinata-go-cli"

# Define the installation directory (where you want to install the binary)
INSTALL_DIR="/usr/local/bin"

# Download the binary
echo "Downloading $BINARY_NAME..."
wget -O "$INSTALL_DIR/$BINARY_NAME" "$BINARY_URL"

# Check if the download was successful
if [ $? -eq 0 ]; then
	echo "Download complete."
	chmod +x "$INSTALL_DIR/$BINARY_NAME" # Make the binary executable
	echo "Setting execute permissions..."
else
	echo "Download failed. Installation aborted."
	exit 1
fi

# Print a success message
echo "Installation complete. You can now use $BINARY_NAME."

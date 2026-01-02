#!/bin/bash
set -e

# Install script for gvid

echo "Installing gvid..."

# Create user if not exists
if ! id -u gvid &>/dev/null; then
    sudo useradd -r -s /bin/false gvid
fi

# Create working directory
sudo mkdir -p /var/lib/gvid
sudo chown gvid:gvid /var/lib/gvid

# Copy binaries
sudo cp bin/gvid /usr/local/bin/
sudo cp bin/gvi-tui /usr/local/bin/
sudo chmod +x /usr/local/bin/gvid /usr/local/bin/gvi-tui

# Install systemd service
sudo cp deploy/gvid.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable gvid

echo "Done. Start with: sudo systemctl start gvid"

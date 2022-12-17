#!/bin/env bash

set -e

app_url="https://github.com/dstoy/cloud-tunnel/raw/master/bin/tunnel"
app_path="/bin/cloud-tunnel"

service_url="https://raw.githubusercontent.com/dstoy/cloud-tunnel/master/systemd.service"
service_file="/etc/systemd/system/cloud-tunnel.service"
service_name="cloud-tunnel"

config_file="/etc/cloud-tunnel.yaml"

# check for systemctl
if ! systemctl --version &>/dev/null; then
    echo "Unsupported system. Could not locate systemd installation"
    exit 1
fi

# download the app
if [ ! -e "$app_path" ]; then
    echo " -> Downloading the application from: $app_url"
    curl -s --show-error -o "$app_path" "$app_url" || {
        echo "Error saving the application to: $app_path"
        exit 1
    }

    echo " -> Saving application to: $app_path"
    chmod +x $app_path
fi

# copy the systemd service file
if [ ! -e "$service_file" ]; then
    echo " -> Configuring systemd service"
    curl -s --show-error -o "$service_file" "$service_url" || {
        echo "Error saving the service to: $service_file"
        exit 1
    }

    echo " -> Enable service: $service_name"
    systemctl enable cloud-tunnel
fi

# create a default configuration
if [ ! -e "$config_file" ]; then
    echo " -> Creating default configuraion"
    tee "$config_file" >/dev/null <<EOF
# queue configuration
queue:
  url: cloud-tunnel
  region: us-east-1

# list of triggers to map events to commands which will be executed
triggers:
  - event: echo "triggered event"
EOF
fi

echo ""
echo "Installation complete"
echo ""
echo "Please edit the configuration file located at: $config_file"
echo "and start the service using the command below:"
echo ""
echo "    systemctl start cloud-tunnel"
echo ""

#!/usr/bin/env bash
# ==============================================================================
# AWS EC2 One-Click Server Initialization Script for Prepio
# Tested on: Ubuntu 22.04 / 24.04 LTS (t2.micro / t3.micro)
# ==============================================================================

set -euo pipefail

echo "=========================================================="
echo "    Prepio AWS EC2 Automated Environment Setup           "
echo "=========================================================="

# 1. Update system packages
echo "--> Updating system packages..."
sudo apt-get update -y
sudo apt-get upgrade -y

# 2. Configure 2GB Swap Space
# EC2 Free Tier (t2.micro/t3.micro) has 1GB RAM. Swap prevents Out-Of-Memory kills.
if [ ! -f /swapfile ]; then
    echo "--> Configuring 2GB swap space for memory stability..."
    sudo fallocate -l 2G /swapfile
    sudo chmod 600 /swapfile
    sudo mkswap /swapfile
    sudo swapon /swapfile
    echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
    sudo sysctl vm.swappiness=10
    echo 'vm.swappiness=10' | sudo tee -a /etc/sysctl.conf
    echo "--> Swap space successfully enabled!"
else
    echo "--> Swap file already exists. Skipping."
fi

# 3. Install Docker and Docker Compose Plugin
echo "--> Installing Docker and dependencies..."
sudo apt-get install -y ca-certificates curl gnupg lsb-release git

sudo install -m 0755 -d /etc/apt/keyrings
if [ ! -f /etc/apt/keyrings/docker.gpg ]; then
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    sudo chmod a+r /etc/apt/keyrings/docker.gpg
fi

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 4. Enable Docker and configure user permissions
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker "$USER" || true

echo "=========================================================="
echo " Setup complete! Docker & Swap are ready."
echo " If you just logged in, please run: newgrp docker"
echo "=========================================================="

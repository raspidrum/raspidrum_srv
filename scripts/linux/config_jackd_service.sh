#!/bin/bash

sudo mkdir -p /etc/systemd/system/jack.service.d/
sudo tee /etc/systemd/system/jack.service.d/10-limits.conf << EOF
[Service]
LimitRTPRIO=95
LimitMEMLOCK=infinity
EOF

echo "session required pam_limits.so" | sudo tee -a /etc/pam.d/login

sudo systemctl daemon-reload

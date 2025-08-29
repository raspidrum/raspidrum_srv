#!/bin/bash

sudo mkdir -p /etc/systemd/system/jack.service.d/
sudo tee /etc/systemd/system/jack.service.d/10-limits.conf << EOF
[Service]
LimitRTPRIO=95
LimitMEMLOCK=infinity
EOF

sudo systemctl daemon-reload

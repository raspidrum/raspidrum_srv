#!/bin/bash

systemctl stop jack.service

export RDRUM_CONFIG=test
cd /opt/raspidrum
./raspidrum

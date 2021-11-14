#!/bin/sh

cd ~/
git clone https://github.com/petabyt/eggnog
cd eggnog

mkdir file

SETUP_USER=$(id -un)

sudo echo "[Unit]\n\
Description=biblesearch service\n\
After=network.target\n\
StartLimitIntervalSec=0\n\
\n\
[Service]\n\
Type=simple\n\
Restart=always\n\
RestartSec=1\n\
User=$(echo $SETUP_USER)\n\
ExecStart=sh -c \"cd ~/eggnog; go run .\"\n\
\n\
[Install]\n\
WantedBy=multi-user.target\n" > /etc/systemd/system/eggnog.service

sudo systemctl start eggnog

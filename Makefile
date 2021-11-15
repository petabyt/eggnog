all:
	go run main.go

init:
	rm -rf file
	mkdir file

systemd:
	sudo sh -c 'echo "[Unit]\n\
	Description=eggnog hosting service\n\
	After=network.target\n\
	StartLimitIntervalSec=0\n\
	\n\
	[Service]\n\
	Type=simple\n\
	Restart=always\n\
	RestartSec=1\n\
	User=$(USER)\n\
	ExecStart=sh -c \"cd ~/eggnog; go run .\"\n\
	\n\
	[Install]\n\
	WantedBy=multi-user.target\n" > /etc/systemd/system/eggnog.service'

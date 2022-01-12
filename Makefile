all:
	go run main.go

init:
	rm -rf file
	mkdir file

systemd:
	echo \
	"[Unit]\n\
	Description=Eggnog image hosting\n\
	After=network.target\n\
	StartLimitIntervalSec=0\n\
	\n\
	[Service]\n\
	Type=simple\n\
	Restart=always\n\
	RestartSec=1\n\
	User=$(USER)\n\
	ExecStart=sh -c 'cd ~/eggnog; go run .'\n\
	\n\
	[Install]\n\
	WantedBy=multi-user.target\n" > eggnog.service

.PHONY: all init systemd

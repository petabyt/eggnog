all:
	go run main.go

init:
	rm -rf file
	mkdir file
	echo "0" > counter

all: clean build run
clean:
	rm -rf ./main
build:
	go build -o main app/main.go
run:
	./main
seed:
	chmod +x curl.sh;./curl.sh
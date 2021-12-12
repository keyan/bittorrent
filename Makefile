run:
	go build && ./bittorrent
test:
	go test ./...
clean:
	find . -name *.out -or -name *.aux -or -name *.log -or -name .*.swp -or -name .*.swo -or -name .DS_Store | xargs -n 1 rm
	rm -f ./bittorrent

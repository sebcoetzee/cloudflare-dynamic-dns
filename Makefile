build:
	env GOOS=darwin GOARCH=amd64 go build -o bin/cloudflare-dynamic-dns-darwin-amd64
	env GOOS=linux GOARCH=amd64 go build -o bin/cloudflare-dynamic-dns-linux-amd64
	env GOOS=linux GOARCH=arm64 go build -o bin/cloudflare-dynamic-dns-linux-arm64
	env GOOS=windows GOARCH=amd64 go build -o bin/cloudflare-dynamic-dns-amd64.exe

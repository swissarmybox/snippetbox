.PHONY: cert run fmt test cover coverdetail buildrun

cert:
	cd tls && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

run:
	go run ./cmd/web

fmt:
	gofmt -w .

test:
	go test ./...

cover:
	go test -cover ./...

coverdetail:
	go test -coverprofile=/tmp/profile.out ./...
	go tool cover -func=/tmp/profile.out

buildrun:
	go build -o /tmp/web ./cmd/web/
	cp -r ./tls /tmp/
	cd /tmp/ && ./web

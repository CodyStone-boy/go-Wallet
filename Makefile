
###############################################################################
# Initial
###############################################################################
goget:
	go get -u -d -v ./...


###############################################################################
# Build
###############################################################################
bld:
	go build -o wallet ./cmd/wallet/main.go
	go build -o coldwallet ./cmd/coldwallet/main.go

run: bld
	./wallet -f 1

###############################################################################
# Test
###############################################################################
gotest:
	go test -v ./...


###############################################################################
# Docker and compose
###############################################################################
bld-docker-go:
	docker build --no-cache -t cayenne-wallet-go:1.10.3 -f ./docker/golang/Dockerfile .


###############################################################################
# Utility
###############################################################################
.PHONY: clean
clean:
	rm -rf detect
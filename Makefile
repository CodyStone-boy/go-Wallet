
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

wallet: bld
	./wallet -f 1

cold-wallet: bld
	./coldwallet -f 4


###############################################################################
# Run
###############################################################################
# TODO:定期的に実行して、動作を確認すること(これを自動化しておきたい)

# 入金データを集約し、未署名のトランザクションを作成する
create-unsigned: bld
	./wallet -f 11

# 未署名のトランザクションに署名する
sign: bld
	./coldwallet -f 5 -i ./data/tx/receipt/receipt_8_unsigned_1534832793024491932

# 署名済トランザクションを送信する
send: bld
	./wallet -f 13 -i ./data/tx/receipt/receipt_8_signed_1534832879778945174

# 出金データから出金トランザクションを作成する
create-payment: bld
	./wallet -f 14

# 出金用に未署名のトランザクションに署名する
sign-payment: bld
	./coldwallet -f 5 -i ./data/tx/payment/payment_3_unsigned_1534832966995082772

# 出金用に署名済トランザクションを送信する
send: bld
	./wallet -f 13 -i ./data/tx/payment/payment_3_signed_1534833088943126101


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
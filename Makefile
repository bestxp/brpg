gen-events:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go_opt=Mpkg/events.proto=github.com/bestxp/brpg/pkg \
		pkg/events.proto


run-client:
	go run cmd/client/*.go

run-server:
	go run cmd/server/*.go
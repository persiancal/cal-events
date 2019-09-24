export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))

test:
	go test ./...

thetool: test
	go build ./cmd/thetool

generate: thetool
	mkdir -p $(ROOT)/dist
	./thetool -file $(ROOT)/jalali.yml generate -dist $(ROOT)/dist

validate: thetool
	./thetool -file $(ROOT)/jalali.yml validate
	./thetool -file $(ROOT)/jalali.yml generate -dist $(ROOT)/dist -compare

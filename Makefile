export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO111MODULE=off
test:
	cd $(ROOT)/cmd/thetool &&  go test

$(ROOT)/cmd/thetool/thetool: test
	cd $(ROOT)/cmd/thetool && go build .

generate: $(ROOT)/cmd/thetool/thetool
	mkdir -p $(ROOT)/dist
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali generate -dist $(ROOT)/dist

validate: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali validate
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali validate-links -ignore

reorder: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali reorder -output - > $(ROOT)/jalali.yaml

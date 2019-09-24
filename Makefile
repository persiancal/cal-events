export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO111MODULE=on
test:
	cd $(ROOT)/cmd/thetool &&  go test 

$(ROOT)/cmd/thetool/thetool: test
	cd $(ROOT)/cmd/thetool && go build .

generate: $(ROOT)/cmd/thetool/thetool
	mkdir -p $(ROOT)/dist
	$(ROOT)/cmd/thetool/thetool -file $(ROOT)/jalali.yml generate -dist $(ROOT)/dist

validate: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -file $(ROOT)/jalali.yml validate
	$(ROOT)/cmd/thetool/thetool -file $(ROOT)/jalali.yml generate -dist $(ROOT)/dist -compare

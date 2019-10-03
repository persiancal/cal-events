export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO111MODULE=off
test:
	cd $(ROOT)/cmd/thetool &&  go test

$(ROOT)/cmd/thetool/thetool: test
	cd $(ROOT)/cmd/thetool && go build .

generate-jalali: $(ROOT)/cmd/thetool/thetool
	mkdir -p $(ROOT)/dist
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali generate -dist $(ROOT)/dist

generate-hijri: $(ROOT)/cmd/thetool/thetool
	mkdir -p $(ROOT)/dist
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri generate -dist $(ROOT)/dist

generate: generate-hijri generate-jalali

validate-jalali: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali validate
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali validate-links -ignore

validate-hirir: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri validate
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri validate-links -ignore
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri generate -dist $(ROOT)/dist -compare

validate: validate-hirir validate-jalali

reorder-jalali: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali reorder -output - > $(ROOT)/jalali.yaml

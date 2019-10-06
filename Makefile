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
	# Make sure update something in dist, since travis looks for changes, and if the build only contains new file
	# skips the deploy
	date > $(ROOT)/dist/.build_at

validate-jalali: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali validate
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali validate-links

validate-hirir: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri validate
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri validate-links

validate: validate-hirir validate-jalali

reorder-jalali: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/jalali reorder -output - > $(ROOT)/jalali.yaml

reorder-hijri: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/hijri reorder -output - > $(ROOT)/hijri.yaml

export ROOT:=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export GO111MODULE=off
test:
	cd $(ROOT)/cmd/thetool &&  go test

$(ROOT)/cmd/thetool/thetool: test
	cd $(ROOT)/cmd/thetool && go build .

split: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/gregorian,$(ROOT)/hijri,$(ROOT)/jalali split -base event -dist $(ROOT)/events

generate: $(ROOT)/cmd/thetool/thetool
	mkdir -p $(ROOT)/dist
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/gregorian,$(ROOT)/events/hijri,$(ROOT)/events/jalali generate -dist $(ROOT)/dist
	# Make sure update something in dist, since travis looks for changes, and if the build only contains new file
	# skips the deploy
	date > $(ROOT)/dist/.build_at

validate: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/gregorian,$(ROOT)/events/hijri,$(ROOT)/events/jalali unique
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/gregorian,$(ROOT)/events/hijri,$(ROOT)/events/jalali validate
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/gregorian,$(ROOT)/events/hijri,$(ROOT)/events/jalali validate-links

reorder-jalali: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/jalali reorder -output - > $(ROOT)/jalali.yaml

reorder-hijri: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/hijri reorder -output - > $(ROOT)/hijri.yaml

reorder-gregorian: $(ROOT)/cmd/thetool/thetool
	$(ROOT)/cmd/thetool/thetool -dir $(ROOT)/events/gregorian reorder -output - > $(ROOT)/gregorian.yaml

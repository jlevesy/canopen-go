
VCAN_DEVICE?=vcan0

.PHONY: vcan
vcan:
	sudo ip link add dev $(VCAN_DEVICE) type vcan
	sudo ip link set up $(VCAN_DEVICE)

.PHONY: gentraffic_vcan
gentraffic_vcan:
	cangen $(VCAN_DEVICE)

.PHONY: dump_vcan
dump_vcan:
	candump $(VCAN_DEVICE)

TEST_COUNT?=1
TEST_PACKAGE?=./...

.PHONY: test
test:
	go test -cover -race -count $(TEST_COUNT) $(TEST_PACKAGE)

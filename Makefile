# NOTE: This Makefile is for supporting comparisons with the Rust
# implementation, not really for local development process.
#
# Assumes installed (not checked): Go toolchain, Rust toolchain, pv.
.DEFAULT_GOAL := bin/base100

bin:
	mkdir -p $@

bin/base100: bin
	go build -o $@ ./cmd/base100

bin/base100-rs: bin
	cargo install base100 \
		--root=/tmp/base100 --bin=base100
	cp /tmp/base100/bin/base100 $@

bin/base100-rs-simd: bin # note: requires Rust Nightly toolchain
	RUSTFLAGS="-C target-cpu=native" cargo install base100 \
		--root=/tmp/base100-simd --bin=base100 \
		--features simd
	cp /tmp/base100-simd/bin/base100 $@

all: bin/base100 bin/base100-rs bin/base100-rs-simd

bench:
	cat /dev/urandom | pv --size 1000000000 -S | bin/base100-rs | bin/base100-rs --decode > /dev/null
	cat /dev/urandom | pv --size 1000000000 -S | bin/base100 | bin/base100 --decode > /dev/null
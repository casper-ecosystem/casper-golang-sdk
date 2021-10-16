module github.com/casper-ecosystem/casper-golang-sdk

go 1.16

require (
	github.com/pkg/errors v0.9.1
	github.com/r3labs/sse/v2 v2.3.6
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.10
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b
)

// currently the library doesn't support bigger buffers which is needed for event stream client so
// temporarily we are using a fork that includes that functionality.
replace github.com/r3labs/sse/v2 v2.3.6 => github.com/Venoox/sse/v2 v2.3.7-0.20211011204928-50b6a32ebb6c

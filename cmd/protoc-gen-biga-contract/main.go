// Command protoc-gen-biga-contract turns the constants declared in
// biga/v1/contract.proto into Go or TypeScript source.
//
// Each enum value that carries a text, number or route option becomes one
// constant named after the value: MEDIA_TYPE_PROTOBUF in TypeScript,
// MediaTypeProtobuf in Go. A route also gets Method and Path constants.
//
// The option is chosen with the plugin parameter lang=go or lang=ts.
package main

import (
	"fmt"
	"io"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail(err)
	}
	var req pluginpb.CodeGeneratorRequest
	if err := proto.Unmarshal(in, &req); err != nil {
		fail(err)
	}
	out, err := proto.Marshal(generate(&req))
	if err != nil {
		fail(err)
	}
	if _, err := os.Stdout.Write(out); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "protoc-gen-biga-contract:", err)
	os.Exit(1)
}

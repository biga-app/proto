package main

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func withUnknown(raw []byte) *descriptorpb.EnumValueOptions {
	o := &descriptorpb.EnumValueOptions{}
	o.ProtoReflect().SetUnknown(raw)
	return o
}

func textOption(s string) *descriptorpb.EnumValueOptions {
	raw := protowire.AppendTag(nil, textField, protowire.BytesType)
	return withUnknown(protowire.AppendString(raw, s))
}

func numberOption(n int64) *descriptorpb.EnumValueOptions {
	raw := protowire.AppendTag(nil, numberField, protowire.VarintType)
	return withUnknown(protowire.AppendVarint(raw, uint64(n)))
}

func routeOption(method, path string) *descriptorpb.EnumValueOptions {
	var msg []byte
	msg = protowire.AppendTag(msg, 1, protowire.BytesType)
	msg = protowire.AppendString(msg, method)
	msg = protowire.AppendTag(msg, 2, protowire.BytesType)
	msg = protowire.AppendString(msg, path)
	raw := protowire.AppendTag(nil, routeField, protowire.BytesType)
	return withUnknown(protowire.AppendBytes(raw, msg))
}

func request(param string) *pluginpb.CodeGeneratorRequest {
	value := func(name string, num int32, o *descriptorpb.EnumValueOptions) *descriptorpb.EnumValueDescriptorProto {
		return &descriptorpb.EnumValueDescriptorProto{Name: proto.String(name), Number: proto.Int32(num), Options: o}
	}
	return &pluginpb.CodeGeneratorRequest{
		Parameter:      proto.String(param),
		FileToGenerate: []string{"biga/v1/contract.proto"},
		ProtoFile: []*descriptorpb.FileDescriptorProto{{
			Name:    proto.String("biga/v1/contract.proto"),
			Package: proto.String("biga.v1"),
			Options: &descriptorpb.FileOptions{GoPackage: proto.String("github.com/x/wirepb/biga/v1;bigav1")},
			EnumType: []*descriptorpb.EnumDescriptorProto{
				{Name: proto.String("MediaType"), Value: []*descriptorpb.EnumValueDescriptorProto{
					value("MEDIA_TYPE_UNSPECIFIED", 0, nil),
					value("MEDIA_TYPE_PROTOBUF", 1, textOption("application/x-protobuf")),
				}},
				{Name: proto.String("Limit"), Value: []*descriptorpb.EnumValueDescriptorProto{
					value("LIMIT_UNSPECIFIED", 0, nil),
					value("LIMIT_MAX_PATH_BYTES", 1, numberOption(4096)),
				}},
				{Name: proto.String("Pattern"), Value: []*descriptorpb.EnumValueDescriptorProto{
					value("PATTERN_SESSION_ID", 1, textOption(`^[A-Za-z0-9]$`)),
				}},
				{Name: proto.String("Route"), Value: []*descriptorpb.EnumValueDescriptorProto{
					value("ROUTE_PUT_SESSION", 1, routeOption("PUT", "/api/sessions/{id}")),
				}},
			},
		}},
	}
}

func TestGo(t *testing.T) {
	resp := generate(request("lang=go"))
	if resp.Error != nil {
		t.Fatal(resp.GetError())
	}
	if len(resp.File) != 1 || resp.File[0].GetName() != "biga/v1/contract_consts.go" {
		t.Fatalf("files = %v", resp.File)
	}
	got := strings.Join(strings.Fields(resp.File[0].GetContent()), " ")
	for _, want := range []string{
		"package bigav1 ",
		`const MediaTypeProtobuf = "application/x-protobuf"`,
		"const LimitMaxPathBytes = 4096",
		`const PatternSessionID = "^[A-Za-z0-9]$"`,
		`RoutePutSession = "PUT /api/sessions/{id}"`,
		`RoutePutSessionMethod = "PUT"`,
		`RoutePutSessionPath = "/api/sessions/{id}"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
	if strings.Contains(got, "Unspecified") {
		t.Errorf("values without an option are skipped:\n%s", got)
	}
}

func TestTypeScript(t *testing.T) {
	resp := generate(request("lang=ts"))
	if resp.Error != nil {
		t.Fatal(resp.GetError())
	}
	got := resp.File[0].GetContent()
	for _, want := range []string{
		`export const MEDIA_TYPE_PROTOBUF = "application/x-protobuf";`,
		"export const LIMIT_MAX_PATH_BYTES = 4096;",
		`export const ROUTE_PUT_SESSION_PATH = "/api/sessions/{id}";`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
}

func TestRejectsUnknownLanguage(t *testing.T) {
	if resp := generate(request("lang=rust")); resp.Error == nil {
		t.Error("an unsupported language should be an error")
	}
	if resp := generate(request("")); resp.Error == nil {
		t.Error("a missing language should be an error")
	}
}

func TestRejectsAmbiguousOptions(t *testing.T) {
	raw := protowire.AppendTag(nil, textField, protowire.BytesType)
	raw = protowire.AppendString(raw, "x")
	raw = protowire.AppendTag(raw, numberField, protowire.VarintType)
	raw = protowire.AppendVarint(raw, 1)
	if _, _, err := readOptions(raw); err == nil {
		t.Error("a value cannot be both text and a number")
	}
}

func TestPascal(t *testing.T) {
	for in, want := range map[string]string{
		"MEDIA_TYPE_PROTOBUF":   "MediaTypeProtobuf",
		"PATTERN_SESSION_ID":    "PatternSessionID",
		"ERROR_CODE_NOT_FOUND":  "ErrorCodeNotFound",
		"LIMIT_HASH_HEX_LENGTH": "LimitHashHexLength",
	} {
		if got := pascal(in); got != want {
			t.Errorf("pascal(%q) = %q, want %q", in, got, want)
		}
	}
}

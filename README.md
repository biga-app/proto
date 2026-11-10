# biga wire protocol

The protobuf schema shared by the CLI, the server and the viewer. The server and
the CLI talk `application/x-protobuf`; the live stream is length-delimited
`StreamEvent` messages (`application/vnd.biga.stream+protobuf`).

Generated code is committed in each consumer, so none of them needs this
repository to build:

| Consumer | Output | Command |
| --- | --- | --- |
| `server` | `internal/wirepb` | `make proto` |
| `cli` | `internal/api/apipb` | `make proto` |
| `website` | `packages/types/src/gen` | `pnpm --filter @biga/types proto` |

Each consumer has a `buf.gen.yaml` that reads this repository from `../proto`.
Generating needs [buf](https://buf.build) and Go on your `PATH`.

## Shared constants

`biga/v1/contract.proto` holds every value the CLI, the server and the viewer
must agree on: media types, the access key header, credential prefixes, error
codes, validation limits, the session ID and hash patterns, and every route.
Each is an enum value with its value in an option:

```proto
ERROR_CODE_NOT_FOUND = 2 [(text) = "not_found"];
LIMIT_MAX_TITLE_RUNES = 4 [(number) = 120];
ROUTE_GET_STATE = 15 [(route) = {method: "GET", path: "/api/sessions/{id}/state"}];
```

`protoc-gen-biga-contract` (in `cmd/`) turns them into constants next to the
generated messages: `ErrorCodeNotFound` in Go, `ERROR_CODE_NOT_FOUND` in
TypeScript, and for a route `RouteGetState`, `RouteGetStateMethod` and
`RouteGetStatePath`. Nothing copies a string or a number by hand, so the three
sides cannot drift. What stays out is what one side alone decides: timeouts,
batch sizes, caches and the limits an operator can change, which the server
reports in `Limits`.

```sh
go test ./...
```

```sh
buf lint
buf breaking --against '.git#branch=main'
```

Field numbers are never reused. Add fields; do not renumber or change types.

## Test vectors

`testdata/e2ee-vectors.json` holds content sealed by the CLI for the
[end-to-end encryption](https://biga.n10nce.dev/docs/encryption) scheme: a key and,
for each case, the plaintext, the blob id and the sealed bytes. The CLI's tests
regenerate it (`BIGA_UPDATE_VECTORS=1 go test ./internal/e2ee`), and the viewer's
tests open every case, which is what keeps the Go and WebCrypto implementations
compatible. Each consumer keeps a copy next to its tests.

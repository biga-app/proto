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

```sh
buf lint
buf breaking --against '.git#branch=main'
```

Field numbers are never reused. Add fields; do not renumber or change types.

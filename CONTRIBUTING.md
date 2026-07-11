# Contributing

Thanks for improving termkit-go.

1. Keep packages framework-agnostic: return strings or values; do not take ownership of a host event loop.
2. Add tests for new rendering or animation behavior.
3. Run `go test ./...`, `go vet ./...`, and `vhs demo.tape` when changing the demo.
4. Keep the README and the relevant file under `docs/` aligned with public API changes.

Open an issue before starting a larger renderer or API change so the public surface stays cohesive.

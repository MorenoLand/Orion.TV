# Contributing

Keep changes focused on OrionTV's desktop client, Wails3 integration, and the remote live page boundary.

Before opening a pull request:

```powershell
npm --prefix frontend install
npm --prefix frontend run build
go test ./...
go vet ./...
go mod verify
git diff --check
```

Do not commit `node_modules`, `frontend/dist`, `bin`, installers, or compressed executable output. Do not launch the application as part of automated checks unless the test explicitly requires an authorized runtime check.

Describe user-facing behavior changes and include the exact platform and toolchain used for build verification.

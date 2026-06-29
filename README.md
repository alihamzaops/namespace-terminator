# namespace-terminator

[![test](https://github.com/alihamzaops/namespace-terminator/actions/workflows/test.yml/badge.svg)](https://github.com/alihamzaops/namespace-terminator/actions/workflows/test.yml)

`nst` is a small cross-platform CLI for force-terminating Kubernetes namespaces stuck in `Terminating`.

It can:
- terminate one or more namespaces by name
- target every namespace currently stuck in `Terminating`, optionally filtered by label
- use your current kubeconfig or an explicit kubeconfig/context
- run on Linux, macOS, and Windows

## Why This Exists

Sometimes a namespace gets stuck in `Terminating` because one or more finalizers never clear. `nst` requests deletion when needed, then clears namespace finalizers through the Kubernetes finalize API.

This is a forceful recovery tool. It can orphan resources that still exist behind the namespace boundary, so use it carefully.

## Install

### Homebrew (macOS / Linux)

```bash
brew tap alihamzaops/homebrew-tap
brew install nst
```

### Scoop (Windows)

```powershell
scoop bucket add alihamzaops https://github.com/alihamzaops/scoop-bucket
scoop install nst
```

### Direct Download

Download the archive for your platform from [GitHub Releases](https://github.com/alihamzaops/namespace-terminator/releases) and place the binary on your `PATH`.

### Go install

```bash
go install github.com/alihamzaops/namespace-terminator/cmd/nst@latest
```

## Usage

Terminate specific namespaces:

```bash
nst payments staging
```

Terminate every namespace stuck in `Terminating`:

```bash
nst --all-terminating
```

Filter stuck namespaces by label:

```bash
nst --all-terminating --selector team=payments
```

Preview what would be targeted without making changes:

```bash
nst --all-terminating --dry-run
nst payments staging --dry-run
```

Use a specific kubeconfig or context:

```bash
nst payments --kubeconfig ~/.kube/config --context prod-cluster
```

Skip the confirmation prompt (required in CI — stdin must be a terminal otherwise):

```bash
nst --all-terminating --yes
```

Emit machine-readable output:

```bash
nst payments --output json
```

## Command Reference

```text
nst <namespace> [<namespace> ...]

Flags:
  --all-terminating      Target all namespaces currently stuck in Terminating
  --context string       Kubeconfig context to use
  --dry-run              Print the namespaces that would be terminated
  --kubeconfig path      Path to a kubeconfig file
  -l, --selector string  Label selector to filter namespaces (requires --all-terminating)
  -o, --output string    Output format: text or json (default "text")
  --timeout duration     Wait time after clearing finalizers (default 15s)
  --version              Print the nst version
  -y, --yes              Skip confirmation for bulk operations
```

## Exit Behavior

- exits `0` when every targeted namespace is deleted successfully
- exits non-zero when any namespace fails to terminate or remains pending after the timeout
- exits non-zero for invalid usage or kubeconfig errors

## How It Works

For each target namespace, `nst`:
1. resolves the target list from explicit names or from namespaces already in `Terminating` (filtered by label if `--selector` is set)
2. sends a namespace delete request if deletion has not started yet
3. clears namespace finalizers through the finalize API
4. waits briefly for the namespace to disappear

When multiple namespaces are targeted, all terminations run concurrently.

## Development

Run locally:

```bash
go run ./cmd/nst --help
```

Run tests:

```bash
go test ./...
go vet ./...
```

Pull requests run tests automatically via the [`test.yml`](.github/workflows/test.yml) workflow.

Releases are published with [GoReleaser](https://goreleaser.com/) on version tags (`v*`). Two secrets are needed on the repository:

| Secret | Purpose |
|--------|---------|
| `TAP_GITHUB_TOKEN` | Write access to `alihamzaops/homebrew-tap` |
| `SCOOP_GITHUB_TOKEN` | Write access to `alihamzaops/scoop-bucket` |

## Security Notes

- Do not commit kubeconfigs, tokens, private keys, or `.env` files to this repository.
- Common secret-bearing files are ignored by `.gitignore`.
- GitHub Actions runs a `gitleaks` secret scan on pushes to `main` and on pull requests.
- You can enable the same check locally with `pre-commit install`, using the repo's `.pre-commit-config.yaml`.

Enable local secret scanning:

```bash
pip install pre-commit
pre-commit install
```

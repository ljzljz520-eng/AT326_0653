# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	arcadeledger/cmd/arcade	[no test files]
--- FAIL: TestArcadeScoreSubmissionIsIdempotent (0.01s)
    arcade_test.go:39: expected idempotent score count 1, got 2
FAIL
FAIL	arcadeledger	0.017s
ok  	arcadeledger/internal/catalog	0.008s
ok  	arcadeledger/internal/model	0.001s
ok  	arcadeledger/internal/query	0.007s
ok  	arcadeledger/internal/schedule	0.009s
ok  	arcadeledger/internal/service	0.023s
ok  	arcadeledger/internal/store	0.020s
ok  	arcadeledger/internal/transport	0.008s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/arcade): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/arcade): exit `0`

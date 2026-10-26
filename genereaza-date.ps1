# F06-T03: versiune sqlc fixată, fără instalare globală.
$ErrorActionPreference = 'Stop'
Push-Location $PSScriptRoot
try {
    go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
} finally {
    Pop-Location
}

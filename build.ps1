# Construiește build/bin/pontix.exe (F05-T13): versiunea din internal/versiune/versiune.go
# ajunge în proprietățile fișierului (prin wails.json) și commit-ul curent prin -ldflags.
# Folosire: .\build.ps1
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

$versiune = (Select-String internal/versiune/versiune.go -Pattern 'const Aplicatie = "(.+)"').Matches[0].Groups[1].Value
$inWails = (Get-Content wails.json -Raw | ConvertFrom-Json).info.productVersion
if ($inWails -ne $versiune) {
    throw "wails.json are productVersion $inWails, dar versiune.go are $versiune. Actualizează wails.json."
}

$commit = git rev-parse --short HEAD
if (git status --porcelain) {
    $commit += '-modificat' # build din fișiere necomise
}

wails build -clean -trimpath -ldflags "-X pontix/internal/versiune.Commit=$commit"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

# Proprietățile se citesc prin Windows Shell, aceeași sursă ca fereastra Proprietăți din Explorer.
$exe = (New-Object -ComObject Shell.Application).Namespace((Resolve-Path build/bin).Path).ParseName('pontix.exe')
$inExe = $exe.ExtendedProperty('System.Software.ProductVersion')
$nume = $exe.ExtendedProperty('System.Software.ProductName')
if ($inExe -ne $versiune -or $nume -ne 'Pontix') {
    throw "pontix.exe are produsul '$nume', versiunea '$inExe'; așteptat 'Pontix', $versiune"
}
Write-Output "build/bin/pontix.exe: $nume $inExe, commit $commit"

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

# Partea numerică a versiunii; câmpurile text le arată Explorer, dar .NET nu le citește
# din resursa scrisă de Wails (limba neutră).
$info = [Diagnostics.FileVersionInfo]::GetVersionInfo((Resolve-Path build/bin/pontix.exe).Path)
$inExe = "$($info.FileMajorPart).$($info.FileMinorPart).$($info.FileBuildPart)"
if ($inExe -ne $versiune) {
    throw "pontix.exe are versiunea $inExe, așteptat $versiune"
}
Write-Output "build/bin/pontix.exe: versiunea $inExe, commit $commit"

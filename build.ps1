# Construiește build/bin/pontix.exe (F05-T13): versiunea din internal/versiune/versiune.go
# ajunge în proprietățile fișierului (prin wails.json) și commit-ul curent prin -ldflags.
# Folosire: .\build.ps1; cu aplicația deschisă: .\build.ps1 -NumeExecutabil pontix-verificare.exe
param([string]$NumeExecutabil = 'pontix.exe')

$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

if ([IO.Path]::GetFileName($NumeExecutabil) -ne $NumeExecutabil -or
    $NumeExecutabil.IndexOfAny([IO.Path]::GetInvalidFileNameChars()) -ge 0 -or
    -not $NumeExecutabil.EndsWith('.exe', [StringComparison]::OrdinalIgnoreCase)) {
    throw 'NumeExecutabil trebuie sa fie un nume de fisier .exe, fara cale.'
}
$caleExecutabil = Join-Path (Join-Path $PSScriptRoot 'build/bin') $NumeExecutabil
if (Test-Path -LiteralPath $caleExecutabil) {
    try {
        $probaAcces = [IO.File]::Open($caleExecutabil, [IO.FileMode]::Open, [IO.FileAccess]::ReadWrite, [IO.FileShare]::None)
        $probaAcces.Dispose()
    } catch {
        throw "Executabilul '$NumeExecutabil' este blocat. Alege alt nume cu -NumeExecutabil; aplicatia deschisa nu va fi oprita."
    }
}

$versiune = (Select-String internal/versiune/versiune.go -Pattern 'const Aplicatie = "(.+)"').Matches[0].Groups[1].Value
$inWails = (Get-Content wails.json -Raw | ConvertFrom-Json).info.productVersion
if ($inWails -ne $versiune) {
    throw "wails.json are productVersion $inWails, dar versiune.go are $versiune. Actualizează wails.json."
}

$commit = git rev-parse --short HEAD
if (git status --porcelain) {
    $commit += '-modificat' # build din fișiere necomise
}

# Folderul bin poate conține karios/ cu datele utilizatorului. Nu se curăță.
wails build -o $NumeExecutabil -trimpath -ldflags "-X pontix/internal/versiune.Commit=$commit"
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

# Proprietățile se citesc prin Windows Shell, aceeași sursă ca fereastra Proprietăți din Explorer.
$exe = (New-Object -ComObject Shell.Application).Namespace((Resolve-Path build/bin).Path).ParseName($NumeExecutabil)
$inExe = $exe.ExtendedProperty('System.Software.ProductVersion')
$nume = $exe.ExtendedProperty('System.Software.ProductName')
if ($inExe -ne $versiune -or $nume -ne 'Pontix') {
    throw "$NumeExecutabil are produsul '$nume', versiunea '$inExe'; așteptat 'Pontix', $versiune"
}
Write-Output "build/bin/${NumeExecutabil}: $nume $inExe, commit $commit"

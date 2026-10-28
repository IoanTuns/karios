# F14-T06: Compilare nativa llama.cpp (llama-server.exe) pe Windows
# Foloseste toolchain-ul portabil local (CMake + MinGW GCC) din tools/toolchain/
# Compileaza in paralel folosind toate nucleele procesorului.
#
# Folosire:
#   .\scripts\compileaza-llama.ps1 [-Tag "b4810"] [-CurataBuild]

param(
    [string]$Tag = "b4810",
    [switch]$CurataBuild
)

$ErrorActionPreference = 'Stop'
$radacinaRepo = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
Set-Location -LiteralPath $radacinaRepo

Write-Host "=== Compilare Nativa llama-server pe Windows ===" -ForegroundColor Cyan

# 1. Asigurare toolchain portabil
$toolchainScript = Join-Path $PSScriptRoot 'pregateste-toolchain.ps1'
$toolchain = & $toolchainScript
$directorTools = Join-Path $radacinaRepo 'tools'

# Adaugare toolchain la PATH temporar pentru sesiunea curenta
$env:PATH = "$($toolchain.CmakeBin);$($toolchain.GccBin);$env:PATH"
$env:CC = Join-Path $toolchain.GccBin 'gcc.exe'
$env:CXX = Join-Path $toolchain.GccBin 'g++.exe'

Write-Host "Verificare versiune compilator:" -ForegroundColor Gray
& gcc.exe --version | Select-Object -First 1 | Write-Host -ForegroundColor Gray
& cmake.exe --version | Select-Object -First 1 | Write-Host -ForegroundColor Gray

# 2. Asigurare surse llama.cpp
$llamaSrc = Join-Path $directorTools 'llama.cpp'
if (-not (Test-Path -LiteralPath $llamaSrc)) {
    Write-Host "Clonare surse llama.cpp (tag: $Tag)..." -ForegroundColor Yellow
    git clone --depth 1 --branch $Tag https://github.com/ggerganov/llama.cpp.git $llamaSrc
    if ($LASTEXITCODE -ne 0) {
        Write-Warning "Clonarea tagului $Tag a esuat; clonare master..."
        git clone --depth 1 https://github.com/ggerganov/llama.cpp.git $llamaSrc
    }
} else {
    Write-Host "Sursele llama.cpp exista deja in: $llamaSrc" -ForegroundColor Green
}

# 3. Director de compilare
$buildDir = Join-Path $directorTools 'llama-build'
if ($CurataBuild -and (Test-Path -LiteralPath $buildDir)) {
    Write-Host "Curatare director de build anterior..." -ForegroundColor Yellow
    Remove-Item -LiteralPath $buildDir -Recurse -Force
}

if (-not (Test-Path -LiteralPath $buildDir)) {
    New-Item -ItemType Directory -Path $buildDir -Force | Out-Null
}

# 4. Configurare CMake
Write-Host "Configurare CMake (MinGW Makefiles, Release, AVX2)..." -ForegroundColor Cyan
$cmakeArgs = @(
    '-B', $buildDir,
    '-S', $llamaSrc,
    '-G', 'MinGW Makefiles',
    '-DCMAKE_BUILD_TYPE=Release',
    '-DLLAMA_BUILD_SERVER=ON',
    '-DLLAMA_BUILD_TESTS=OFF',
    '-DLLAMA_BUILD_EXAMPLES=OFF',
    '-DLLAMA_AVX2=ON'
)
& cmake.exe @cmakeArgs
if ($LASTEXITCODE -ne 0) { throw "Configurarea CMake a esuat cu codul $LASTEXITCODE." }

# 5. Compilare paralela a serverului
$nuclee = $env:NUMBER_OF_PROCESSORS
if (-not $nuclee) { $nuclee = 4 }
Write-Host "Compilare llama-server folosind $nuclee nuclee paralele..." -ForegroundColor Cyan

$momentStart = [DateTimeOffset]::Now
$buildArgs = @(
    '--build', $buildDir,
    '--config', 'Release',
    '--target', 'llama-server',
    '-j', $nuclee
)
& cmake.exe @buildArgs
if ($LASTEXITCODE -ne 0) { throw "Compilarea a esuat cu codul $LASTEXITCODE." }

$durata = [DateTimeOffset]::Now - $momentStart
Write-Host ("Compilare finalizata cu succes in: {0:mm}m {0:ss}s" -f $durata) -ForegroundColor Green

# 6. Cautare executabil generat
$posibileCaiBinar = @(
    (Join-Path $buildDir 'bin/llama-server.exe'),
    (Join-Path $buildDir 'bin/Release/llama-server.exe'),
    (Join-Path $buildDir 'llama-server.exe')
)

$binarCompilat = $null
foreach ($cale in $posibileCaiBinar) {
    if (Test-Path -LiteralPath $cale -PathType Leaf) {
        $binarCompilat = $cale
        break
    }
}

if (-not $binarCompilat) {
    $cautare = Get-ChildItem -Path $buildDir -Filter 'llama-server.exe' -Recurse
    if ($cautare.Count -gt 0) {
        $binarCompilat = $cautare[0].FullName
    } else {
        throw "Executabilul llama-server.exe nu a fost gasit in directorul de build."
    }
}

$hashBinar = (Get-FileHash -LiteralPath $binarCompilat -Algorithm SHA256).Hash
$marimeBytes = (Get-Item -LiteralPath $binarCompilat).Length
$marimeMB = [Math]::Round($marimeBytes / 1MB, 2)

Write-Host "Binar generat: $binarCompilat ($marimeMB MB)" -ForegroundColor Green
Write-Host "SHA256: $hashBinar" -ForegroundColor Gray

# 7. Copiere in karios/ai/ si .cache-ai/
$tintaKarios = Join-Path $radacinaRepo 'karios/ai'
if (-not (Test-Path -LiteralPath $tintaKarios)) {
    New-Item -ItemType Directory -Path $tintaKarios -Force | Out-Null
}
Copy-Item -LiteralPath $binarCompilat -Destination (Join-Path $tintaKarios 'llama-server.exe') -Force
Write-Host "Copiat in: $tintaKarios\llama-server.exe" -ForegroundColor Green

$tintaCache = Join-Path $radacinaRepo '.cache-ai'
if (-not (Test-Path -LiteralPath $tintaCache)) {
    New-Item -ItemType Directory -Path $tintaCache -Force | Out-Null
}
Copy-Item -LiteralPath $binarCompilat -Destination (Join-Path $tintaCache 'llama-server.exe') -Force
Write-Host "Copiat in: $tintaCache\llama-server.exe" -ForegroundColor Green

[PSCustomObject]@{
    Executabil = $binarCompilat
    SHA256     = $hashBinar
    MarimeMB   = $marimeMB
    Durata     = $durata
}


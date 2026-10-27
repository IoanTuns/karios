# F14-T06: Pregatire toolchain C/C++ si CMake portabil pe Windows
# Descarca si configureaza local CMake si GCC (w64devkit) in tools/toolchain/
# Fara instalare pe sistem, fara drepturi de administrator.

param(
    [string]$DirectorTools = ""
)

$ErrorActionPreference = 'Stop'
$radacinaRepo = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
if (-not $DirectorTools) {
    $DirectorTools = Join-Path $radacinaRepo 'tools'
}

$toolchainDir = Join-Path $DirectorTools 'toolchain'
$cmakeDir = Join-Path $toolchainDir 'cmake'
$gccDir = Join-Path $toolchainDir 'w64devkit'

if (-not (Test-Path -LiteralPath $toolchainDir)) {
    New-Item -ItemType Directory -Path $toolchainDir -Force | Out-Null
}

# 1. Verificare / Descarcare CMake Portabil
$cmakeExe = $null
$posibilCmake = Get-ChildItem -Path $cmakeDir -Filter 'cmake.exe' -Recurse -ErrorAction SilentlyContinue
if ($posibilCmake.Count -gt 0) {
    $cmakeExe = $posibilCmake[0].FullName
}

if (-not $cmakeExe) {
    Write-Host "Descarcare CMake portabil pentru Windows x64..." -ForegroundColor Cyan
    $urlCmake = 'https://github.com/Kitware/CMake/releases/download/v3.31.5/cmake-3.31.5-windows-x86_64.zip'
    $zipCmake = Join-Path $toolchainDir 'cmake.zip'
    Invoke-WebRequest -Uri $urlCmake -OutFile $zipCmake
    Write-Host "Dezarhivare CMake..." -ForegroundColor Yellow
    Expand-Archive -LiteralPath $zipCmake -DestinationPath $cmakeDir -Force
    Remove-Item -LiteralPath $zipCmake -Force
    $posibilCmake = Get-ChildItem -Path $cmakeDir -Filter 'cmake.exe' -Recurse
    if ($posibilCmake.Count -gt 0) {
        $cmakeExe = $posibilCmake[0].FullName
    } else {
        throw "Nu s-a gasit cmake.exe dupa dezarhivare."
    }
}

# 2. Verificare / Descarcare w64devkit (GCC 14 + G++ + Make)
$gccExe = $null
$posibilGcc = Get-ChildItem -Path $gccDir -Filter 'gcc.exe' -Recurse -ErrorAction SilentlyContinue
if ($posibilGcc.Count -gt 0) {
    $gccExe = $posibilGcc[0].FullName
}

if (-not $gccExe) {
    Write-Host "Descarcare w64devkit (GCC 14 + G++ + Make portabil)..." -ForegroundColor Cyan
    $urlW64 = 'https://github.com/skeeto/w64devkit/releases/download/v2.10.1/w64devkit-x64-2.10.1.7z.exe'
    $sfxW64 = Join-Path $toolchainDir 'w64devkit-setup.exe'
    if (-not (Test-Path -LiteralPath $gccDir)) {
        New-Item -ItemType Directory -Path $gccDir -Force | Out-Null
    }
    Invoke-WebRequest -Uri $urlW64 -OutFile $sfxW64
    Write-Host "Dezarhivare w64devkit..." -ForegroundColor Yellow
    try {
        & tar.exe -xf $sfxW64 -C $gccDir
    } catch {
        & $sfxW64 -y "-o$gccDir"
    }
    Remove-Item -LiteralPath $sfxW64 -Force -ErrorAction SilentlyContinue
    $posibilGcc = Get-ChildItem -Path $gccDir -Filter 'gcc.exe' -Recurse
    if ($posibilGcc.Count -gt 0) {
        $gccExe = $posibilGcc[0].FullName
    } else {
        throw "Nu s-a gasit gcc.exe dupa dezarhivare."
    }
}

$binCmake = [IO.Path]::GetDirectoryName($cmakeExe)
$binGcc = [IO.Path]::GetDirectoryName($gccExe)

Write-Host "Toolchain portabil configurat cu succes:" -ForegroundColor Green
Write-Host "  CMake: $cmakeExe"
Write-Host "  GCC:   $gccExe"

[PSCustomObject]@{
    CmakeExe = $cmakeExe
    CmakeBin = $binCmake
    GccExe   = $gccExe
    GccBin   = $binGcc
}


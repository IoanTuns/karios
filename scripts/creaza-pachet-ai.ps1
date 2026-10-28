# F14-T05: Creare pachet portabil Pontix AI Ready
# Automatizează pregătirea unui director sau a unei arhive portabile all-in-one
# conținând Pontix, motorul llama-server.exe și modelul GGUF în karios/ai/.
#
# Folosire:
#   .\scripts\creaza-pachet-ai.ps1 [-NumeModel "qwen2.5-3b-instruct-q4_k_m.gguf"] [-ArhiveazaZip] [-DescarcaDacaLipseste]

param(
    [string]$NumeModel = 'qwen2.5-3b-instruct-q4_k_m.gguf',
    [string]$UrlModel = 'https://huggingface.co/Qwen/Qwen2.5-3B-Instruct-GGUF/resolve/main/qwen2.5-3b-instruct-q4_k_m.gguf',
    [string]$UrlLlamaServer = 'https://github.com/ggerganov/llama.cpp/releases/download/b4810/llama-b4810-bin-win-avx2-x64.zip',
    [switch]$DescarcaDacaLipseste,
    [switch]$CuCompilareAI,
    [switch]$ArhiveazaZip
)

$ErrorActionPreference = 'Stop'
$radacinaRepo = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
Set-Location -LiteralPath $radacinaRepo

Write-Host "=== Pregatire Pachet Pontix AI Ready ===" -ForegroundColor Cyan

# 1. Determinare versiune din codul sursă
$versiune = (Select-String internal/versiune/versiune.go -Pattern '^const Aplicatie = "(.+)"').Matches[0].Groups[1].Value
if ($versiune -notmatch '^\d+\.\d+\.\d+$') { throw 'Versiunea aplicatiei este invalida.' }
$commit = git rev-parse --short HEAD

# 2. Cautare ultimul build sau compilare nouă
$radacinaBuild = Join-Path $radacinaRepo 'build/bin'
$ultimulBuildExe = $null
if (Test-Path -LiteralPath $radacinaBuild) {
    $executabile = Get-ChildItem -Path $radacinaBuild -Filter 'pontix.exe' -Recurse | Sort-Object LastWriteTime -Descending
    if ($executabile.Count -gt 0) {
        $ultimulBuildExe = $executabile[0].FullName
    }
}

if (-not $ultimulBuildExe) {
    Write-Host "Nu a fost gasit un build existent. Se ruleaza .\build.ps1..." -ForegroundColor Yellow
    $rez = & (Join-Path $radacinaRepo 'build.ps1')
    $ultimulBuildExe = $rez.Executabil
}

Write-Host "Executabil sursa: $ultimulBuildExe" -ForegroundColor Green

# 3. Directorul pachetului AI Ready
$moment = [DateTimeOffset]::Now
$guidPachet = [Guid]::NewGuid().ToString('N')
$numePachet = "pontix-$versiune-ai-ready-" + $moment.ToString('yyyyMMdd-HHmmss') + "-$guidPachet"
$destinatiePachet = Join-Path $radacinaBuild $numePachet
New-Item -ItemType Directory -Path $destinatiePachet | Out-Null

# Copiere executabil principal
$destExe = Join-Path $destinatiePachet 'pontix.exe'
Copy-Item -LiteralPath $ultimulBuildExe -Destination $destExe

# 4. Creare structura karios/ai/
$destKariosAI = Join-Path $destinatiePachet 'karios/ai'
New-Item -ItemType Directory -Path $destKariosAI -Force | Out-Null

# Verificare surse AI: cache local .cache-ai/ sau folder local karios/ai/
$folderCache = Join-Path $radacinaRepo '.cache-ai'
$llamaGasit = $null
$modelGasit = $null

$posibileLocatiiLlama = @(
    (Join-Path $folderCache 'llama-server.exe'),
    (Join-Path $radacinaRepo 'karios/ai/llama-server.exe')
)
foreach ($loc in $posibileLocatiiLlama) {
    if (Test-Path -LiteralPath $loc -PathType Leaf) {
        $llamaGasit = $loc
        break
    }
}

$posibileLocatiiModel = @(
    (Join-Path $folderCache $NumeModel),
    (Join-Path $radacinaRepo "karios/ai/$NumeModel")
)
foreach ($loc in $posibileLocatiiModel) {
    if (Test-Path -LiteralPath $loc -PathType Leaf) {
        $modelGasit = $loc
        break
    }
}

# 5. Compilare sau descărcare dacă lipsește și a fost cerut explicit
if ($CuCompilareAI -or (-not $llamaGasit -and -not $DescarcaDacaLipseste)) {
    Write-Host "Compilare nativa llama-server via scripts/compileaza-llama.ps1..." -ForegroundColor Cyan
    $compScript = Join-Path $PSScriptRoot 'compileaza-llama.ps1'
    if (Test-Path -LiteralPath $compScript) {
        $rezComp = & $compScript
        $llamaGasit = $rezComp.Executabil
    }
}

if ($DescarcaDacaLipseste) {
    if (-not (Test-Path -LiteralPath $folderCache)) {
        New-Item -ItemType Directory -Path $folderCache | Out-Null
    }

    if (-not $llamaGasit) {
        Write-Host "Descarcare motor llama-server de la $UrlLlamaServer..." -ForegroundColor Yellow
        $zipLlama = Join-Path $folderCache 'llama-bin.zip'
        Invoke-WebRequest -Uri $UrlLlamaServer -OutFile $zipLlama
        Expand-Archive -LiteralPath $zipLlama -DestinationPath (Join-Path $folderCache 'llama-extrase') -Force
        $gasite = Get-ChildItem -Path (Join-Path $folderCache 'llama-extrase') -Filter 'llama-server.exe' -Recurse
        if ($gasite.Count -gt 0) {
            Copy-Item -LiteralPath $gasite[0].FullName -Destination (Join-Path $folderCache 'llama-server.exe')
            $llamaGasit = Join-Path $folderCache 'llama-server.exe'
        }
    }

    if (-not $modelGasit) {
        Write-Host "Descarcare model $NumeModel (aprox. 2 GB)..." -ForegroundColor Yellow
        $destModel = Join-Path $folderCache $NumeModel
        Invoke-WebRequest -Uri $UrlModel -OutFile $destModel
        if (Test-Path -LiteralPath $destModel) {
            $modelGasit = $destModel
        }
    }
}

# 6. Copiere fișiere AI în pachet dacă sunt disponibile
$componenteAI = [ordered]@{}
if ($llamaGasit) {
    Write-Host "Copiere llama-server.exe in pachet..." -ForegroundColor Green
    Copy-Item -LiteralPath $llamaGasit -Destination (Join-Path $destKariosAI 'llama-server.exe')
    $componenteAI['llama_server'] = (Get-FileHash -LiteralPath $llamaGasit -Algorithm SHA256).Hash
} else {
    Write-Host "Atentie: llama-server.exe nu a fost gasit in cache." -ForegroundColor DarkYellow
}

if ($modelGasit) {
    Write-Host "Copiere model $NumeModel in pachet..." -ForegroundColor Green
    Copy-Item -LiteralPath $modelGasit -Destination (Join-Path $destKariosAI $NumeModel)
    $componenteAI['model'] = $NumeModel
    $componenteAI['sha256_model'] = (Get-FileHash -LiteralPath $modelGasit -Algorithm SHA256).Hash
} else {
    Write-Host "Atentie: Modelul $NumeModel nu a fost gasit in cache." -ForegroundColor DarkYellow
}

# 7. Script utilitar pentru descărcare la cerere în destinație
$scriptDescarcareContinut = @"
# Script pentru descarcarea modelului AI in folderul portabil karios/ai/
`$ErrorActionPreference = 'Stop'
`$destFolder = `$PSScriptRoot
Write-Host "Descarcare model local Qwen 2.5 3B GGUF in `$destFolder..." -ForegroundColor Cyan
`$modelUrl = "$UrlModel"
`$tinta = Join-Path `$destFolder "$NumeModel"
Invoke-WebRequest -Uri `$modelUrl -OutFile `$tinta
Write-Host "Modelul a fost descarcat cu succes! Porniti Pontix pentru a folosi Asistentul AI." -ForegroundColor Green
"@
[IO.File]::WriteAllText((Join-Path $destKariosAI 'descarca-model.ps1'), $scriptDescarcareContinut, [Text.UTF8Encoding]::new($false))

# 8. Creare README-AI-READY.txt
$readmeContinut = @"
========================================================================
PONTIX $versiune - EDITIE PORTABILA AI READY
========================================================================

Aceasta versiune contine aplicatia Pontix configurata pentru utilizare
impreuna cu asistentul de limbaj natural LLM local.

CONFIDENTIALITATE SI SECURITATE STRICTA (Decizia D042):
- Motorul AI ruleaza 100% offline, pe calculatorul dumneavoastra (127.0.0.1).
- Nicio data cu privire la angajati, contracte sau pontaj nu este trimisa
  vreodata in cloud sau pe servere externe.
- Modelul AI propune doar modificari validate, pe care le confirmati
  dumneavoastra explicit inainte de salvare.

STRUCTURA:
- pontix.exe             : Aplicatia principala (dublu-click pentru rulare).
- karios/ai/             : Contine motorul local (llama-server.exe) si
                           modelul optimizat de limbaj ($NumeModel).

UTILIZARE RAPIDA:
1. Rulati pontix.exe.
2. Navigati in meniul lateral la "Asistent AI".
3. Daca motorul nu este pornit automat, apasati "Porneste motorul".
4. Puteti scrie comenzi in limba romana (ex: "Pune-o pe Popescu de garda vineri").

Dezvoltat de Ioan Tuns. Toate drepturile rezervate (c) 2026.
"@
[IO.File]::WriteAllText((Join-Path $destinatiePachet 'README-AI-READY.txt'), $readmeContinut, [Text.UTF8Encoding]::new($false))

# 9. Scriere manifest build.json pentru pachet
$manifest = [ordered]@{
    tip_pachet = 'ai_ready'
    versiune = $versiune
    commit = $commit
    creat_la = $moment.ToString('o')
    executabil = 'pontix.exe'
    sha256_executabil = (Get-FileHash -LiteralPath $destExe -Algorithm SHA256).Hash
    componente_ai = $componenteAI
}
$caleManifest = Join-Path $destinatiePachet 'build.json'
[IO.File]::WriteAllText($caleManifest, (($manifest | ConvertTo-Json) -replace "`r`n", "`n") + "`n", [Text.UTF8Encoding]::new($false))

# 10. Arhivare opțională ZIP
$caleArhiva = $null
if ($ArhiveazaZip) {
    Write-Host "Generare arhiva ZIP: $numePachet.zip..." -ForegroundColor Yellow
    $caleArhiva = Join-Path $radacinaBuild "$numePachet.zip"
    Compress-Archive -Path "$destinatiePachet\*" -DestinationPath $caleArhiva -Force
    Write-Host "Arhiva ZIP generata: $caleArhiva" -ForegroundColor Green
}

Write-Host "`n=== Pachet AI Ready Creat cu Succes ===" -ForegroundColor Cyan
Write-Host "Director: $destinatiePachet" -ForegroundColor White
if ($caleArhiva) {
    Write-Host "Arhiva:   $caleArhiva" -ForegroundColor White
}
Write-Host "SHA256 pontix.exe: $($manifest.sha256_executabil)" -ForegroundColor Gray

[PSCustomObject]@{
    Director = $destinatiePachet
    Arhiva = $caleArhiva
    Versiune = $versiune
    Commit = $commit
}


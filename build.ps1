# F02-T12: fiecare build are propriul director; datele existente nu se suprascriu.
# Folosire: .\build.ps1 [-NumeExecutabil pontix.exe] [-CuCompilareAI]
param([string]$NumeExecutabil = 'pontix.exe', [switch]$CuCompilareAI)

$ErrorActionPreference = 'Stop'
$locatieInitiala = Get-Location
$temporarInitial = $env:TEMP
$tmpInitial = $env:TMP
$directorBuild = $null
$directorTemporar = $null
$manifest = $null

try {
    Set-Location -LiteralPath $PSScriptRoot
    if ([string]::IsNullOrWhiteSpace($NumeExecutabil) -or
        [IO.Path]::GetFileName($NumeExecutabil) -ne $NumeExecutabil -or
        $NumeExecutabil.IndexOfAny([IO.Path]::GetInvalidFileNameChars()) -ge 0 -or
        -not $NumeExecutabil.EndsWith('.exe', [StringComparison]::OrdinalIgnoreCase)) {
        throw 'NumeExecutabil trebuie sa fie un nume de fisier .exe, fara cale.'
    }

    $versiune = (Select-String internal/versiune/versiune.go -Pattern '^const Aplicatie = "(.+)"').Matches[0].Groups[1].Value
    if ($versiune -notmatch '^\d+\.\d+\.\d+$') { throw 'Versiunea aplicatiei este invalida.' }
    $inWails = (Get-Content -LiteralPath wails.json -Raw | ConvertFrom-Json).info.productVersion
    if ($inWails -ne $versiune) { throw "wails.json are productVersion $inWails, dar versiune.go are $versiune." }
    $commit = git rev-parse --short HEAD
    if ($LASTEXITCODE -ne 0) { throw 'Commitul Git nu poate fi citit.' }
    $modificari = git status --porcelain
    if ($LASTEXITCODE -ne 0) { throw 'Starea Git nu poate fi citita.' }
    if ($modificari) { $commit += '-modificat' }

    $radacinaIesire = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot 'build/bin'))
    $identificator = [Guid]::NewGuid().ToString('N')
    $momentBuild = [DateTimeOffset]::Now
    $generatorManual = Join-Path $PSScriptRoot 'docs/manual/genereaza.py'
    # Manualul trebuie revizuit pentru versiunea curenta inaintea compilarii costisitoare.
    python $generatorManual --versiune $versiune --commit $commit --data $momentBuild.ToString('o') --verifica-sursa
    if ($LASTEXITCODE -ne 0) { throw 'Verificarea sursei manualului a esuat.' }
    $numeDirector = 'pontix-' + $versiune + '-' + $momentBuild.ToString('yyyyMMdd-HHmmss') + '-' + $identificator
    $directorBuild = [IO.Path]::GetFullPath((Join-Path $radacinaIesire $numeDirector))
    if ([IO.Path]::GetDirectoryName($directorBuild) -ne $radacinaIesire -or
        (Test-Path -LiteralPath $directorBuild)) { throw 'Directorul de build nu este nou sau nu apartine build/bin.' }
    New-Item -ItemType Directory -Path $directorBuild | Out-Null
    $directorTemporar = Join-Path $directorBuild '.temporar'
    New-Item -ItemType Directory -Path $directorTemporar | Out-Null
    $env:TEMP = $directorTemporar
    $env:TMP = $directorTemporar

    $manifest = [ordered]@{
        versiune = $versiune
        commit = $commit
        creat_la = $momentBuild.ToString('o')
        executabil = $NumeExecutabil
        stare = 'in_lucru'
    }
    $caleManifest = Join-Path $directorBuild 'build.json'
    function Scrie-ManifestBuild {
        [IO.File]::WriteAllText($caleManifest, (($manifest | ConvertTo-Json -Depth 8) -replace "`r`n", "`n") + "`n", [Text.UTF8Encoding]::new($false))
    }
    Scrie-ManifestBuild
    # Motorul AI compilat nativ se face inaintea aplicatiei: suma lui intra in executabil, care
    # porneste doar servere cunoscute (D056).
    $flagi = "-X pontix/internal/versiune.Commit=$commit"
    $rezAI = $null
    if ($CuCompilareAI) {
        Write-Host "Compilare nativa llama-server pentru livrare..." -ForegroundColor Cyan
        $rezAI = & (Join-Path $PSScriptRoot 'scripts/compileaza-llama.ps1')
        if ($rezAI.SHA256 -notmatch '^[0-9A-Fa-f]{64}$') { throw 'Suma llama-server compilat este invalida.' }
        $flagi += " -X pontix/internal/ai.HashServerCompilat=$($rezAI.SHA256)"
    }
    # -m/-nosyncgomod pastreaza dependentele fixate. Nu folosim -clean sau nivelul verbose 2.
    $iesireWails = $numeDirector + '/' + $NumeExecutabil
    wails build -o $iesireWails -m -nosyncgomod -trimpath -ldflags $flagi | Out-Host
    $codBuild = $LASTEXITCODE
    # Generatorul adauga spatii la sfarsit de linie; pastram formatul fisierului urmarit in Git.
    $caleModele = Join-Path $PSScriptRoot 'frontend/wailsjs/go/models.ts'
    if (Test-Path -LiteralPath $caleModele) {
        $modeleInitiale = [IO.File]::ReadAllText($caleModele)
        $modeleNormalizate = [regex]::Replace($modeleInitiale.Replace([string][char]13,''),'[ \t]+$','',[Text.RegularExpressions.RegexOptions]::Multiline)
        if ($modeleNormalizate -cne $modeleInitiale) { [IO.File]::WriteAllText($caleModele,$modeleNormalizate,[Text.UTF8Encoding]::new($false)) }
    }
    if ($codBuild -ne 0) { throw "Buildul Wails a esuat cu codul $codBuild. Director: $directorBuild" }

    $caleExecutabil = Join-Path $directorBuild $NumeExecutabil
    if (-not (Test-Path -LiteralPath $caleExecutabil -PathType Leaf)) { throw 'Executabilul nu a fost produs in directorul buildului.' }
    # Aceleasi proprietati Windows ca in fereastra Explorer.
    $exe = (New-Object -ComObject Shell.Application).Namespace($directorBuild).ParseName($NumeExecutabil)
    $inExe = $exe.ExtendedProperty('System.Software.ProductVersion')
    $nume = $exe.ExtendedProperty('System.Software.ProductName')
    if ($inExe -ne $versiune -or $nume -ne 'Pontix') { throw "Executabilul are produsul '$nume', versiunea '$inExe'; asteptat Pontix $versiune." }
    $manifest.sha256 = (Get-FileHash -LiteralPath $caleExecutabil -Algorithm SHA256).Hash

    if ($null -ne $rezAI) {
        $destAiDir = Join-Path $directorBuild 'karios/ai'
        New-Item -ItemType Directory -Path $destAiDir -Force | Out-Null
        Copy-Item -LiteralPath $rezAI.Executabil -Destination (Join-Path $destAiDir 'llama-server.exe') -Force
        $manifest.llama_server_sha256 = $rezAI.SHA256
    }

    # O livrare nu este completa fara manualele aceleiasi versiuni si aceluiasi commit.
    $jsonManuale = python $generatorManual --versiune $versiune --commit $commit --data $momentBuild.ToString('o') --destinatie $directorBuild --temporare $directorTemporar
    if ($LASTEXITCODE -ne 0) { throw 'Generarea manualelor a esuat; livrarea ramane incompleta.' }
    $manifest.manuale = ($jsonManuale -join "`n") | ConvertFrom-Json
    if ($manifest.manuale.versiune -ne $versiune -or $manifest.manuale.commit -ne $commit -or $manifest.manuale.fisiere.Count -ne 3) {
        throw 'Manualele nu corespund livrarii.'
    }
    foreach ($manual in $manifest.manuale.fisiere) {
        $caleManual = Join-Path $directorBuild $manual.nume
        if ([IO.Path]::GetFileName($manual.nume) -ne $manual.nume -or
            -not (Test-Path -LiteralPath $caleManual -PathType Leaf) -or
            (Get-FileHash -LiteralPath $caleManual -Algorithm SHA256).Hash -ne $manual.sha256) {
            throw 'Fisierul manualului lipseste sau are suma de verificare incorecta.'
        }
    }
    $manifest.stare = 'complet'
    Scrie-ManifestBuild
    # F02-T22: livrarea contine numai fisierele aprobate, niciodata baza/seiful.
    $calePachet = python (Join-Path $PSScriptRoot 'verificari/pachet_release.py') --director $directorBuild
    if ($LASTEXITCODE -ne 0) { throw 'Generarea pachetului ZIP a esuat.' }
    Write-Host $calePachet
    Write-Host "$caleExecutabil`: Pontix $versiune, commit $commit"
    # Rezultatul poate fi preluat de verificari sau de un script de livrare.
    [PSCustomObject]@{Director = $directorBuild; Executabil = $caleExecutabil; Manifest = $caleManifest; Pachet = $calePachet}
} catch {
    if ($null -ne $manifest) {
        $manifest.stare = 'esuat'
        Scrie-ManifestBuild
    }
    throw
} finally {
    $env:TEMP = $temporarInitial
    $env:TMP = $tmpInitial
    Set-Location -LiteralPath $locatieInitiala
    # Eliminam numai temporarele create in aceasta invocare, niciodata datele buildurilor existente.
    if ($null -ne $directorTemporar -and (Test-Path -LiteralPath $directorTemporar)) {
        $tintaTemporare = [IO.Path]::GetFullPath($directorTemporar)
        if ([IO.Path]::GetDirectoryName($tintaTemporare) -ne $directorBuild -or
            [IO.Path]::GetFileName($tintaTemporare) -ne '.temporar') { throw 'Tinta temporarelor nu apartine buildului curent.' }
        $legaturi = @(Get-ChildItem -LiteralPath $tintaTemporare -Force -Recurse | Where-Object { $_.Attributes -band [IO.FileAttributes]::ReparsePoint })
        if (((Get-Item -LiteralPath $tintaTemporare -Force).Attributes -band [IO.FileAttributes]::ReparsePoint) -or $legaturi.Count -gt 0) {
            Write-Warning "Temporarele contin o legatura si au fost pastrate: $tintaTemporare"
        } else {
            try { Remove-Item -LiteralPath $tintaTemporare -Recurse -Force }
            catch { Write-Warning "Unele temporare sunt blocate si au fost pastrate in $tintaTemporare" }
        }
    }
}

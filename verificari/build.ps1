# F02-T12: doua builduri izolate, temporare, manifest si esec; compilatorul este simulat.
# ExecutabilReferinta trebuie sa fie un build real al versiunii curente.
param([Parameter(Mandatory=$true)][string]$ExecutabilReferinta)
$ErrorActionPreference = 'Stop'
$radacinaProiect = Split-Path $PSScriptRoot -Parent
$scriptBuild = Join-Path $radacinaProiect 'build.ps1'
$folderProbe = [IO.Path]::GetFullPath((Join-Path $radacinaProiect 'build/bin'))
$referinta = (Resolve-Path -LiteralPath $ExecutabilReferinta).Path
$apeluri = [Collections.Generic.List[string]]::new()
$simuleazaEsec = $false
$temporarInitial = $env:TEMP
$tmpInitial = $env:TMP
$locatieInitiala = (Get-Location).Path
$numeProba = 'proba-conservare-' + [Guid]::NewGuid().ToString('N') + '.exe'
$caleProba = Join-Path $folderProbe $numeProba
$blocare = $null

function wails {
    $indice = [Array]::IndexOf([object[]]$args, '-o')
    if ($indice -lt 0 -or $args -contains '-clean') { throw 'Argumente Wails nesigure.' }
    $cale = [IO.Path]::GetFullPath((Join-Path $folderProbe $args[$indice+1]))
    $director = [IO.Path]::GetDirectoryName($cale)
    if ([IO.Path]::GetDirectoryName($director) -ne $folderProbe) { throw 'Buildul nu are subdirector propriu.' }
    $apeluri.Add($director)
    if ($env:TEMP -ne (Join-Path $director '.temporar') -or $env:TMP -ne $env:TEMP) { throw 'Temporarele nu sunt izolate.' }
    $dateTemporare = Join-Path $env:TEMP 'karios'
    New-Item -ItemType Directory -Path $dateTemporare | Out-Null
    [IO.File]::WriteAllText((Join-Path $dateTemporare 'proba.txt'), 'numai date fictive ale generatorului')
    if ($simuleazaEsec) { Set-Variable -Name LASTEXITCODE -Value 17 -Scope 1; return }
    Copy-Item -LiteralPath $referinta -Destination $cale
    Set-Variable -Name LASTEXITCODE -Value 0 -Scope 1
}

try {
    foreach ($numeInvalid in @('../alt.exe', 'sub/alt.exe', 'pontix.txt', '', 'C:\alt.exe')) {
        $refuzat = $false
        try { & $scriptBuild -NumeExecutabil $numeInvalid | Out-Null } catch { $refuzat = $_.Exception.Message -like '*fara cale*' }
        if (-not $refuzat) { throw "Numele invalid a fost acceptat: '$numeInvalid'" }
    }
    if ($apeluri.Count -ne 0) { throw 'Un nume invalid a invocat compilatorul.' }
    [IO.File]::WriteAllText($caleProba, 'fisier existent protejat')
    $blocare = [IO.File]::Open($caleProba, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::None)
    $primul = & $scriptBuild -NumeExecutabil $numeProba
    $alDoilea = & $scriptBuild -NumeExecutabil $numeProba
    if ($primul.Director -eq $alDoilea.Director) { throw 'Buildurile au acelasi director.' }
    foreach ($rezultat in @($primul,$alDoilea)) {
        $manifest = Get-Content -LiteralPath $rezultat.Manifest -Raw | ConvertFrom-Json
        if ($manifest.stare -ne 'complet' -or $manifest.sha256 -ne (Get-FileHash -LiteralPath $rezultat.Executabil).Hash) { throw 'Manifestul nu descrie executabilul verificat.' }
        if (Test-Path -LiteralPath (Join-Path $rezultat.Director '.temporar')) { throw 'Temporarele buildului nu au fost curatate.' }
    }
    $simuleazaEsec = $true
    $refuzat = $false
    try { & $scriptBuild -NumeExecutabil $numeProba | Out-Null } catch { $refuzat = $_.Exception.Message -like '*codul 17*' }
    if (-not $refuzat) { throw 'Esec Wails nesemnalat.' }
    $esuat = $apeluri[$apeluri.Count-1]
    $manifestEsuat = Get-Content -LiteralPath (Join-Path $esuat 'build.json') -Raw | ConvertFrom-Json
    if ($manifestEsuat.stare -ne 'esuat' -or (Test-Path -LiteralPath (Join-Path $esuat '.temporar'))) { throw 'Esec fara manifest sau temporare ramase.' }
    if ($env:TEMP -ne $temporarInitial -or $env:TMP -ne $tmpInitial -or (Get-Location).Path -ne $locatieInitiala) { throw 'Mediul apelantului nu a fost restaurat.' }
    $blocare.Dispose()
    $blocare = $null
    if ([IO.File]::ReadAllText($caleProba) -ne 'fisier existent protejat') { throw 'Fisierul existent a fost modificat.' }
    Write-Output 'Build verificat: directoare distincte, iesire existenta blocata pastrata, hash/manifest, temporare izolate, cai invalide si esec Wails.'
} finally {
    if ($null -ne $blocare) { $blocare.Dispose() }
    if (Test-Path -LiteralPath $caleProba) { Remove-Item -LiteralPath $caleProba }
    foreach ($director in $apeluri) {
        $tinta = [IO.Path]::GetFullPath($director)
        if ([IO.Path]::GetDirectoryName($tinta) -ne $folderProbe -or [IO.Path]::GetFileName($tinta) -notmatch '^pontix-\d+\.\d+\.\d+-\d{8}-\d{6}-[a-f0-9]{32}$') { throw 'Director de test neasteptat; nu se sterge.' }
        Remove-Item -LiteralPath $tinta -Recurse -Force
    }
}

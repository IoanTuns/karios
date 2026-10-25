# F05-T15: verifică refuzul căilor și al executabilului blocat, fără build real.
$ErrorActionPreference = 'Stop'
$radacinaProiect = Split-Path $PSScriptRoot -Parent
$scriptBuild = Join-Path $radacinaProiect 'build.ps1'

foreach ($numeInvalid in @('../alt.exe', 'sub/alt.exe', 'pontix.txt', '')) {
    $refuzat = $false
    try { & $scriptBuild -NumeExecutabil $numeInvalid } catch { $refuzat = $_.Exception.Message -like '*fara cale*' }
    if (-not $refuzat) { throw "Numele invalid a fost acceptat: '$numeInvalid'" }
}

$folderProbe = [IO.Path]::GetFullPath((Join-Path $radacinaProiect 'build/bin'))
if (-not (Test-Path -LiteralPath $folderProbe)) { New-Item -ItemType Directory -Path $folderProbe -Force | Out-Null }
$numeProba = 'proba-blocare-' + [Guid]::NewGuid().ToString('N') + '.exe'
$caleProba = [IO.Path]::GetFullPath((Join-Path $folderProbe $numeProba))
if ([IO.Path]::GetDirectoryName($caleProba) -ne $folderProbe) { throw 'Calea probei nu este în build/bin.' }
[IO.File]::WriteAllText($caleProba, 'proba exclusiva F05-T15')
$blocare = $null
try {
    $blocare = [IO.File]::Open($caleProba, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::None)
    $refuzat = $false
    try { & $scriptBuild -NumeExecutabil $numeProba } catch { $refuzat = $_.Exception.Message -like '*este blocat*' }
    if (-not $refuzat) { throw 'Executabilul blocat nu a fost refuzat înainte de build.' }
} finally {
    if ($null -ne $blocare) { $blocare.Dispose() }
    # Singurul fișier eliminat este proba creată de acest script, verificată mai sus.
    Remove-Item -LiteralPath $caleProba
}
Write-Output 'Build verificat: cai invalide si fisier blocat refuzate inainte de compilare; datele aplicatiei nu sunt accesate.'

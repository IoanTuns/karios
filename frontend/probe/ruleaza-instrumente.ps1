param([int]$Port = 5177)
$ErrorActionPreference = 'Stop'
$radacinaProba = Split-Path (Split-Path $PSScriptRoot -Parent) -Parent
$directorProba = Join-Path $radacinaProba 'build/bin/proba-ui-F14-T25'
$null = New-Item -ItemType Directory -Force -Path $directorProba
$browserProba = 'C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe'
$esecuriProba = 0
$cazuriProba = @('clarificari', 'sabloane', 'plan', 'anulare', 'conflict', 'expirat', 'regresie-editare', 'regresie-istoric', 'regresie-rezultat', 'regresie-eroare', 'regresie-previzualizare', 'regresie-fara-motor')
foreach ($limbaProba in @('ro', 'en')) {
    foreach ($cazProba in $cazuriProba) {
        $numeProba = "$cazProba-$limbaProba"
        $profilProba = Join-Path $directorProba ('profil-' + [guid]::NewGuid().ToString('N'))
        $domProba = Join-Path $directorProba "$numeProba.html"
        $logProba = Join-Path $directorProba "$numeProba.log"
        $imagineProba = Join-Path $directorProba "$numeProba.png"
        $latimeProba = if ($limbaProba -eq 'en') { 800 } else { 1024 }
        $paginaProba = if ($cazProba.StartsWith('regresie-')) { 'asistent-ui' } else { 'instrumente-ai' }
        $cazBrowserProba = $cazProba -replace '^regresie-', ''
        $argumenteProba = @('--headless=new', '--disable-gpu', '--no-sandbox', '--no-first-run', '--disable-extensions', "--user-data-dir=`"$profilProba`"", '--virtual-time-budget=6000', "--window-size=$latimeProba,733", '--dump-dom', "--screenshot=`"$imagineProba`"", "http://127.0.0.1:$Port/probe/$paginaProba.html?caz=$cazBrowserProba&limba=$limbaProba")
        $procesProba = Start-Process -FilePath $browserProba -ArgumentList $argumenteProba -WindowStyle Hidden -PassThru -RedirectStandardOutput $domProba -RedirectStandardError $logProba
        if (!$procesProba.WaitForExit(45000)) { $procesProba.Kill(); $esecuriProba++; Write-Output "$numeProba : browser expirat"; continue }
        $continutProba = ''
        for ($incercareProba = 0; $incercareProba -lt 40; $incercareProba++) {
            try {
                $fluxProba = [IO.File]::Open($domProba, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::ReadWrite)
                $cititorProba = New-Object IO.StreamReader($fluxProba)
                try { $continutProba = $cititorProba.ReadToEnd() } finally { $cititorProba.Dispose() }
            } catch [IO.IOException] { Start-Sleep -Milliseconds 100; continue }
            if ($continutProba -match '</html>') { break }
            Start-Sleep -Milliseconds 100
        }
        if ($continutProba -match 'data-rezultat="trecut"') { Write-Output "$numeProba : trecut" }
        else { $esecuriProba++; $motivProba = [regex]::Match($continutProba, 'data-rezultat="([^"]*)"').Groups[1].Value; Write-Output "$numeProba : EȘEC $motivProba" }
    }
}
if ($esecuriProba -gt 0) { throw "$esecuriProba probe eșuate; DOM și capturi în $directorProba" }

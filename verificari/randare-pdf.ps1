# Verificare vizuală offline cu motorul PDF inclus în Windows; fără instalări.
param([Parameter(Mandatory=$true)][string]$PDF,[Parameter(Mandatory=$true)][string]$PNG,[int]$Pagina=1)
$ErrorActionPreference='Stop'
Add-Type -AssemblyName System.Runtime.WindowsRuntime
$null=[Windows.Storage.StorageFile,Windows.Storage,ContentType=WindowsRuntime]
$null=[Windows.Data.Pdf.PdfDocument,Windows.Data.Pdf,ContentType=WindowsRuntime]
$null=[Windows.Storage.Streams.InMemoryRandomAccessStream,Windows.Storage.Streams,ContentType=WindowsRuntime]
$null=[Windows.Storage.Streams.DataReader,Windows.Storage.Streams,ContentType=WindowsRuntime]
$null=[Windows.Data.Pdf.PdfPageRenderOptions,Windows.Data.Pdf,ContentType=WindowsRuntime]
function Asteapta-Rezultat($Operatie,[Type]$Tip) {
    $metoda=[System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object { $_.Name -eq 'AsTask' -and $_.IsGenericMethodDefinition -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncOperation`1' } | Select-Object -First 1
    $sarcina=$metoda.MakeGenericMethod($Tip).Invoke($null,@($Operatie))
    $sarcina.Wait()
    return $sarcina.Result
}
$fisier=Asteapta-Rezultat ([Windows.Storage.StorageFile]::GetFileFromPathAsync((Resolve-Path -LiteralPath $PDF).Path)) ([Windows.Storage.StorageFile])
$document=Asteapta-Rezultat ([Windows.Data.Pdf.PdfDocument]::LoadFromFileAsync($fisier)) ([Windows.Data.Pdf.PdfDocument])
if ($Pagina -lt 1 -or $Pagina -gt $document.PageCount) { throw 'Pagina inexistenta.' }
$foaie=$document.GetPage($Pagina-1)
$memorie=[Windows.Storage.Streams.InMemoryRandomAccessStream]::new()
try {
    $optiuni=[Windows.Data.Pdf.PdfPageRenderOptions]::new()
    $optiuni.DestinationWidth=2000
    $actiune=$foaie.RenderToStreamAsync($memorie,$optiuni)
    $metodaActiune=[System.WindowsRuntimeSystemExtensions].GetMethods() | Where-Object { $_.Name -eq 'AsTask' -and -not $_.IsGenericMethod -and $_.GetParameters().Count -eq 1 -and $_.GetParameters()[0].ParameterType.Name -eq 'IAsyncAction' } | Select-Object -First 1
    $sarcinaRandare=$metodaActiune.Invoke($null,@($actiune))
    $sarcinaRandare.Wait()
    $cititor=[Windows.Storage.Streams.DataReader]::new($memorie.GetInputStreamAt(0))
    try {
        $null=Asteapta-Rezultat ($cititor.LoadAsync([uint32]$memorie.Size)) ([uint32])
        $octeti=New-Object byte[] ([int]$memorie.Size)
        $cititor.ReadBytes($octeti)
        [IO.File]::WriteAllBytes([IO.Path]::GetFullPath($PNG),$octeti)
    } finally { $cititor.Dispose() }
} finally { $memorie.Dispose();$foaie.Dispose() }
Write-Output $PNG

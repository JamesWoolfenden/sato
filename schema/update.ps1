Add-Type -AssemblyName System.IO.Compression.FileSystem
function Unzip {
    param(
        [Parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]
        [string]$zipfile,

        [Parameter(Mandatory=$true)]
        [ValidateNotNullOrEmpty()]
        [string]$outpath
    )

    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $outpath)
}

$root = "./"| Resolve-Path
$schemaUrl = "https://schema.cloudformation.us-east-1.amazonaws.com/CloudformationSchema.zip"
$filepath = Join-Path $root "CloudformationSchema.zip"
write-host "path $filepath"
Get-ChildItem *.json| ForEach-Object { Remove-Item $_}

try {
    Write-Progress -Activity "Downloading Schema" -Status "Downloading..."
    invoke-webrequest $schemaUrl -OutFile $filepath
    Write-Progress -Activity "Downloading Schema" -Completed
    Unzip $filepath $root
} catch {
    Write-Error "Failed to download/extract schema: $_"
    exit 1
}

Remove-Item $filepath

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

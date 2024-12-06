param(
    [Parameter(Mandatory = $false)]
    [ValidateNotNullOrEmpty()]
    [string]$message = "new release"
)

$versionPattern = '^\d+\.\d+\.\d+$'
$version = $null

try
{
    $version = $( git describe --tags --abbrev=0 ) -replace "v"
    if ($version -notmatch $versionPattern)
    {
        Write-Error "Invalid version format $version. Expected: x.y.z"
        exit 1
    }

    $splitter = $version.split(".")
    $build = [int]($splitter[2]) + 1
    [string]$newVersion = $splitter[0] + "." + $splitter[1] + "." + $build.ToString()

    if ([version]$newVersion -le [version]$version)
    {
        Write-Error "New version must be greater than current version"
        exit 1
    }

    Write-Host "Current version: $version"
    Write-Host "New version: $newVersion"
    Write-Host "Creating new tag..."

    git tag -a v$newVersion -m "$message"
    git push origin v$newVersion
}
catch
{
    Write-Error "An error occurred: $_"
    exit 1
}

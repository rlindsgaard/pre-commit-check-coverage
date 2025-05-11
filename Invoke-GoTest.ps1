<#
    .SYNOPSIS
    Command line script to invoke tests with flags to produce
    coverage reports.

    .DESCRIPTION
    Invokes 'go test' for all packages in the current module and
    ensures that a coverage report is enabled and then processed.
#>
[CmdletBinding()]
param()

Push-Location
$exitCode = $null

$exit = {
    Write-Verbose "Exiting"
    Pop-Location
    Write-Verbose "Unregistering event"
    Unregister-Event -SourceIdentifier Powershell.Exiting
    if($null -eq $exitCode) {
        Write-Error "Abnormal exit"
        exit 1
    }
    Write-Verbose "Exiting normally"
    Write-Debug "Exiting with code: $($exitCode)"
    exit $exitCode
}

Register-EngineEvent -SourceIdentifier Powershell.Exiting -Action $exit | Out-Null

Set-Location $PSScriptRoot

Write-Verbose "Invoking tests"
go test ./... -coverprofile .\cover.out
$exitCode=$LASTEXITCODE
Write-Debug "exitCode is now: $($exitCode)"
Write-Verbose "Processing coverage file"
go run .\helpers\extract_relative_paths.go -moduleName 'github.com/rlindsgaard/pre-commit-check-coverage' -root $PWD.Path -coverfile .\cover.out | Out-File 'sha256sums.txt'

Invoke-Command -ScriptBlock $exit
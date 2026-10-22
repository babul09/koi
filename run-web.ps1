$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

$goExe = (Get-Command go -ErrorAction SilentlyContinue).Source
if (-not $goExe) {
    $goCandidate = Join-Path $env:ProgramFiles 'Go\bin\go.exe'
    if (Test-Path $goCandidate) { $goExe = $goCandidate }
    else { throw 'Go was not found. Install Go and reopen PowerShell.' }
}

$env:GOOS = 'js'
$env:GOARCH = 'wasm'
& $goExe build -o game.wasm .
$motionBuildExit = $LASTEXITCODE
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
if ($motionBuildExit -ne 0) { throw 'WebAssembly build failed.' }

$goroot = & $goExe env GOROOT
Copy-Item (Join-Path $goroot 'lib\wasm\wasm_exec.js') '.\wasm_exec.js' -Force
Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Write-Host 'Motion DNA is serving at http://localhost:8000 (Ctrl+C to stop).'
& $goExe run .\cmd\serve

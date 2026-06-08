Write-Host "Building U-Hermes..." -ForegroundColor Green

# Build frontend
Write-Host "[1/2] Building React frontend..."
Set-Location web
npm install
npm run build
Set-Location ..

# Build Go binary with embedded frontend
Write-Host "[2/2] Building Go binary..."
$env:CGO_ENABLED = "1"
go build -ldflags "-s -w -H windowsgui" -o u-hermes.exe .

Write-Host "Build complete: u-hermes.exe" -ForegroundColor Green

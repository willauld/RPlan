
param(
    [string]$file = ""
)
if ($file -eq "") {
    echo 'Need to define $file'
    exit
}
echo "go run ARetirementPlanner.go tomldata.go $file -AMSdlmo"
go run ARetirementPlanner.go tomldata.go $file -AMSdlmo

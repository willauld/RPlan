
param(
    [string]$file = "",
    [string]$compare = $false 
)
if ($file -eq "") {
    echo 'Need to define $file'
    exit
}
$base = [io.path]::GetFileNameWithoutExtension("$file")
$flags = "-AMSkmo"
$temp = "$base" 
$logName1 = "$temp.log"
$csvName1 = "$temp.csv"
echo "will do>> go run ARetirementPlanner.go tomldata.go $file --logfile=$logName1 --csv=$csvName1 $flags"

go run ARetirementPlanner.go tomldata.go $file --logfile=$logName1 --csv=$csvName1 $flags -d
#go run ARetirementPlanner.go tomldata.go $file $flags > $logName1

if ($compare -eq $true) {
    $temp = "$base" + "RO"
    $logName2 = "$temp.log"
    $csvName2 = "$temp.csv"
    #.\ro.exe $file --logfile=$logName2 $flags
    echo ".\ro.exe $file --logfile=$logName2 --csv=$csvName2 $flags"
    .\ro.exe $file --logfile=$logName2 --csv=$csvName2 $flags
    #dobash "diff $logName1 $logName2 "
}

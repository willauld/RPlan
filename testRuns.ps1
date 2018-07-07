
param(
    [int]$stop_count = 3
)
#echo "running testRun.ps1"
echo "Stop Count is: $stop_count"

cd .\tempResults

$tests_strmap = Get-ChildItem ".\" -Filter *.strmap 
#$tests_toml = Get-ChildItem ".\" -Filter *.toml
#$all_tests = $test_strmap + $test_toml

#$flags = "-AkvdMmo"
$flags = "-AkvdM"

#echo $tests_strmap
$All_OK = $true 
for ($i = 0; $i -lt $tests_strmap.Count; $i++) {
    $test_file = $tests_strmap[$i].FullName
    #echo $test_file
    $test_basefile = $tests_strmap[$i].BaseName
    echo "Processing: $test_basefile"
    #cd .\tempResults
    #pwd

    cp $test_basefile".log" $test_basefile".knowngood.log"
    cp $test_basefile".csv" $test_basefile".knowngood.csv"

    go run  ..\ARetirementPlanner.go  ..\tomldata.go $test_file --logfile=$test_basefile".log" --csv=$test_basefile".csv" $flags

    compare-object (get-content $test_basefile".log") (get-content $test_basefile".knowngood.log") > $test_basefile".diff"

    #fc.exe $test_basefile".log" $test_basefile".knowngood.log" > $test_basefile".diff"
    
    dobash "diff $test_basefile.log $test_basefile.knowngood.log > $test_basefile.diff"

    if (Test-Path $test_basefile".diff") { 
        if ((Get-Item $test_basefile".diff").length -gt 0kb) {
            echo "  Has difference"
            $All_OK = $false
        }
        else {
            echo "  matches known good"
            rm $test_basefile".diff"
        }
    }
    if ($i -eq $stop_count) {
        echo "IT IS ZERO"
        break
    }
    #cd ..
}
if ($All_OK) {
    echo "*** All is OK ***"
}
else {
    echo "*** Something did not match known good logs ***"
}
cd $PSScriptRoot

#Foreach-Object {
#    $content = Get-Content $_.FullNam
#
#    #filter and save content to the original file
#    $content | Where-Object {$_ -match 'step[49]'} | Set-Content $_.FullName
#
#    #filter and save content to a new file 
#    $content | Where-Object {$_ -match 'step[49]'} | Set-Content ($_.BaseName + '_out.log')
##}
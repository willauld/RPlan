param(
    [string]$vermajor = "0", 
    [string]$verminor = "4",
    [string]$verStr = "alpha",
    [switch]$resetPatch = $false,
    [switch]$tagRepo = $false
)

# file has current patch number mananged by this script
$verPatchFile = ".patchNum" 

#Test-Path $verPatchFile -PathType Leaf
if ($resetPatch -or !(Test-Path $verPatchFile -PathType Leaf)) {
    $zero = 0
    Set-Content $verPatchFile -value $zero 
}
$verpatch = Get-Content .patchNum
[int]$verpatchInt = $verpatch
Write-Output "verpatch is: $verpatch, Int $verpatchInt"

$newPatchNum = 1 + $verpatchInt
[string]$newPatchStr = $newPatchNum

Set-Content $verPatchFile -value $newPatchStr
$verpatch = Get-Content .patchNum
[int]$verpatchInt = $verpatch
Write-Output "verpatch is: $verpatch, Int $verpatchInt"

if ($tagRepo) {
    $tagStr = "v$vermajor.$verminor.$verpatch"
    git tag -a $tagStr -m "Tag version $tagStr"
    git push origin $tagStr

    # to delete tags if needed:
    #git tag --delete $tagStr
    #git push --delete origin $tagStr
}

$libgitver = git -C C:\home\auld\godev\src\github.com\willauld\rplanlib describe --always

Write-Output "rplanlib git hash: $libgitver"

$liblpsimplexgitver = git -C C:\home\auld\godev\src\github.com\willauld\lpsimplex describe --always

Write-Output "lpsimplex git hash: $liblpsimplexgitver"

$drivergitver = git describe --always --tags --long

Write-Output "ro driver git hash: $drivergitver"

$a = Get-Date
$btime = $a.ToUniversalTime().ToString() -replace ' ', '_'

go build -ldflags "-X main.vermajor=$vermajor -X main.verminor=$verminor -X main.verpatch=$verpatch -X main.verstr=$verStr -X main.gitDriverHash=$drivergitver -X main.gitLibHash=$libgitver -X main.buildTime=$btime -X main.gitlpsimplexHash=$liblpsimplexgitver " -o ro.exe -v


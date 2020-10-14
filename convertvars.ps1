param(
    
    [Parameter(Mandatory=$false)]
    $env,
    [Parameter(Mandatory=$false)]
    $imagetag,
    [Parameter(Mandatory=$false)]
    $registry
    [Parameter(Mandatory=$false)]
    $storageaccount
)
$variables = @{}
if($env){
    $variables | Add-Member -MemberType noteproperty -Name env -Value $env -Force
}
if($imagetag){
    $variables | Add-Member -MemberType noteproperty -Name imagetag -Value $imagetag -Force
}
if($registry){
    $variables | Add-Member -MemberType noteproperty -Name registry -Value $registry -Force
}
if($storageaccount){
    $variables | Add-Member -MemberType noteproperty -Name storageaccount -Value $storageaccount -Force
}
$variables | ConvertTo-Json | Out-File variables.json
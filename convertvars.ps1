param(
    
    [Parameter(Mandatory)]
    $env,
    [Parameter(Mandatory)]
    $imagetag,
    [Parameter(Mandatory)]
    $registry
)
$variables = @{}
$variables | Add-Member -MemberType noteproperty -Name env -Value $env -Force
$variables | Add-Member -MemberType noteproperty -Name imagetag -Value $imagetag -Force
$variables | Add-Member -MemberType noteproperty -Name registry -Value $registry -Force

$variables | ConvertTo-Json | Out-File variables.json
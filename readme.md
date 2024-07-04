GROP
This is a learning project for go (clone of my rust attemp), creating a custom version of GREP for use on windows, or linux

Features

search text highlighting
glob checking (*.txt or *.json)
case senstive (or insenitive) -i


format grop.exe [options] <query> <file>

supports recursive searching with -r / --recursive
supports wild card for <file> like: *.json or *.go
has intergrated tests >> go test 
    Testing the string searching

go build
cp .\grop.exe c:\utils\grop.exe
edit path variable and add C:\utilis\ as a folder restart windows (if already added restart terminal)

TODO: 
    Add Line number
    

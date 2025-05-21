# GROP
This is a learning project for go, creating a custom version of GREP for use on windows, or linux

## Features

search text highlighting
glob checking (*.txt or *.json)
case senstive (or insenitive) -i
Add Line number


format grop.exe [options] <query> <file>

Usage: grop [options] <query> <file>
Search for patterns using the <query> in each file

Options:
    -i, --ignore-case   query will ignore casing (default will include character cases)
    -f, --filename      query searching will print matches with filename
    -r, --recursive     query will recursivly search all subdirectories
    -l, --linenumber    matches will indicate linenumer that they were found at

supports recursive searching with -r / --recursive
supports wild card for <file> like: *.json or *.go
has intergrated tests >> go test 
    Testing the string searching

## Building

go build
cp .\grop.exe c:\utils\grop.exe
edit path variable and add C:\utilis\ as a folder restart windows (if already added restart terminal)

## Testing

go test ./... -v




## TODO: 
    
    -c --count print the number of watches inside the file (this will include filename)

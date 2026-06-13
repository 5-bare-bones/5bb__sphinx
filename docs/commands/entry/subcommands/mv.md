## Use

`sphinx file mv <src> <dst>`

## Description

Move an entry.
TODO: change docs to fit entry

In case any of the paths contains spaces within it, it must be enclosed by double quotes.

## Flags

No flags.

## Examples

Move a file:
```
sphinx file mv oldFile newFile
```

Move a directory:
```
sphinx file mv oldDir/ newDir/
```

Move a file into a directory:
```
sphinx file mv oldDir/test.txt newDir/
```

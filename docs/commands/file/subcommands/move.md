## Use

`sphinx file move {src} {dst}`

*Aliases*: `mv`.ß

## Description

Move a file or directory.

In case any of the paths contains spaces within it, it must be enclosed by double quotes.

## Flags

No flags.

## Examples

Move a file:
```
sphinx file move oldFile newFile
```

Move a directory:
```
sphinx file move oldDir/ newDir/
```
# TODO: maybe we want a directory command?

Move a file into a directory:
```
sphinx file move oldDir/test.txt newDir/
```

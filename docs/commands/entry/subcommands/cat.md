## Use

`sphinx file cat <name> [-c copy]`

## Description

Read entry and write to standard output

TODO: change docs to fit entry OR deleteß

## Flags 

|  Name     | Shorthand |     Type      |    Default    |              Description              |
|-----------|-----------|---------------|---------------|---------------------------------------|
| copy      | c         | bool          | false         | Copy file content to the clipboard    |

### Examples

Write one file:
```
sphinx cat fileName
```

Write one file and copy content to the clipboard:
```
sphinx cat fileName -c
```

Write multiple files:
```
sphinx cat file1 file2 file3
```

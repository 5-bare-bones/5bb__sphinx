## Use

`sphinx file show <name> [-c copy]`

## Description

Read file and write to standard output

## Flags 

|  Name     | Shorthand |     Type      |    Default    |              Description              |
|-----------|-----------|---------------|---------------|---------------------------------------|
| copy      | c         | bool          | false         | Copy file content to the clipboard    |

### Examples

Write one file:
```
sphinx show fileName
```

Write one file and copy content to the clipboard:
```
sphinx show fileName -c
```

Write multiple files:
```
sphinx show file1 file2 file3
```

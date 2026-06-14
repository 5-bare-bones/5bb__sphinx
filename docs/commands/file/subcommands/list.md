## Use

`sphinx file list {name} [-f filter]`

## Description

List files.

## Flags

|  Name     | Shorthand |     Type      |    Default    |      Description      |
|-----------|-----------|---------------|---------------|-----------------------|
| filter    | f         | bool          | false         | Filter files by name  |

### Example

List trip.txt file and copy its content to the clipboard:
```
sphinx file list trip.txt
```

Filter among files:
```
sphinx file list book -f
```

List all the files:
```
sphinx file list
```

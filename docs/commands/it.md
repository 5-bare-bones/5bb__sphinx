## Use

`sphinx it <command|flags|name>`

## Description

Interactive prompt.		
This command behaves depending on the arguments received, it requests the missing information.

|       Given       |    Requests       |
|-------------------|-------------------|
| command           | flags and name    |
| command and flags | name              |
| name              | command and flags |

## Flags 

No flags.

### Examples

No arguments:
```
sphinx it
```

Command without flags:
```
sphinx it list
```

Command with flags:
```
sphinx it list -s -q
```

Only the name:
```
sphinx sample
```

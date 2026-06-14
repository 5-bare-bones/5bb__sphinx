## Use

`sphinx prompt {command|flags|name}`

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
sphinx prompt
```

Command without flags:
```
sphinx prompt list
```

Command with flags:
```
sphinx prompt list -s -q
```

Only the name:
```
sphinx sample
```

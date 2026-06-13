## Use 

`sphinx card list <name> [-f filter] [-q qr] [-s show]`

## Description

List cards.

## Flags

|  Name     | Shorthand |     Type      |    Default    |                 Description                   |
|-----------|-----------|---------------|---------------|-----------------------------------------------|
| filter    | f         | bool          | false         | Filter cards                                  |
| qr        | q         | bool          | false         | Display the number QR code on the terminal   	|
| show      | s         | bool          | false         | Show card number and security code            |

### Examples

List a card showing sensitive information:
```
sphinx card list Sample -s
```

Filter:
```
sphinx file list Sample -f
```

List all cards;
```
sphinx card list
```

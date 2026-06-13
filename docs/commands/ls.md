## Use

`sphinx ls <name> [-f filter] [-q qr] [-s show]`

*Aliases*: entries, list.

## Description

List entries.

> Listing all the entries does not check for expired entries, this decision was taken to prevent high loads when the number of entries is elevated. Listing a single entry does notifies if it is expired.

## Flags 

|  Name     | Shorthand |     Type      |    Default    |                                  Description                                         	|
|-----------|-----------|---------------|---------------|---------------------------------------------------------------------------------------|
| filter    | f         | bool          | false         | Filter entries                                                                       	|
| qr        | q         | bool          | false         | Display the password QR code on the terminal (not-available when listing all entries)	|
| show      | s         | bool          | false         | Show entry password                                                                  	|

### Examples

List an entry:
```
sphinx ls Sample
```

List one and show sensitive information:
```
sphinx ls Sample -s
```

Filter:
```
sphinx ls Sample -f
```

List all entries:
```
sphinx ls
```

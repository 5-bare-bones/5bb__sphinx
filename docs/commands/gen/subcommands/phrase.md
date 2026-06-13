## Use

`sphinx gen phrase [-c copy] [-l length] [-s separator] [-i include] [-e exclude] [-m mute] [-L list] [-q qr]`

*Aliases*: passphrase.

## Description

Generate a random passphrase.

## Flags

|  Name     | Shorthand |     Type      |    Default    |                         Description                           |
|-----------|-----------|---------------|---------------|---------------------------------------------------------------|
| copy      | c         | bool          | false         | Copy the passphrase to the clipboard                          |
| length    | l         | uint64        | 0             | Passphrase length                                             |
| separator | s         | string        | " " (space)   | Character that separates each word                            |
| include   | i         | []string      | nil           | Words to include in the passphrase                            |
| exclude   | e         | []string      | nil           | Words to exclude from the passphrase                          |
| list      | L         | string        | "WordList"    | Passphrase generating method (NoList, WordList, SyllableList) |
| qr        | q         | bool          | false         | Display the passphrase QR code image on the terminal			|
| mute      | m         | bool          | false         | Mute standard output when the passphrase is copied            |

### Examples

Generate a passphrase without a list (default):
```
sphinx gen phrase --length 6
```

Generate a passphrase with word list:
```
sphinx gen phrase --length 7 --list WordList
```

Generate a passphrase with syllable list:
```
sphinx gen phrase --length 12 --list SyllableList
```

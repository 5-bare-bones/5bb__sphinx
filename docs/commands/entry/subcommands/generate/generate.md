## Use

`sphinx generate [-c|--copy] [-l|--length length] [-L|--levels levels] [-i|--include include] [-e|--exclude exclude] [-m|--mute mute] [-r|--repeat repeat] [-q|--qr]`

## Description

Generate a random password.

### Subcommands

- `sphinx generate phrase`: Generate a random passphrase.

## Flags

|  Name     | Shorthand |     Type      |    Default    |                   Description                     |
|-----------|-----------|---------------|---------------|---------------------------------------------------|
| copy      | c         | bool          | false         | Create an entry with a custom password            |
| length    | l         | uint64        | 0             | Password length                                   |
| levels    | L         | []int         | [1,2,3,4,5]   | Password levels                                   |
| include   | i         | string        | ""            | Characters to include in the password             |
| exclude   | e         | string        | ""            | Characters to exclude from the password           |
| repeat    | r         | bool          | true          | Character repetition                              |
| qr        | q         | bool          | false         | Display the password QR code on the terminal		|
| mute      | m         | bool          | false         | Mute standard output when the password is copied 	|

### Format levels

> Default is [1, 2, 3, 4, 5].

1. Lowercase characters (a, b, c...)
2. Uppercase characters (A, B, C...)
3. Digits (0, 1, 2...)
4. Space
5. Special characters (!, $, %...)

### Examples

Generate a password:
```
sphinx generate phrase --levels 1,2,3,4,5 --length 16 --include s4^%$
```

Generate and show the QR code image:
```
sphinx generate phrase --length 20 --qr
```

Generate, copy and mute standard output:
```
sphinx generate phrase --length 25 --copy --mute
```

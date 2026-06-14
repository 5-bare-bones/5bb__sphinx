## Use

`sphinx topt {name} [-c copy] [-i info] [-t timeout]`

## Description

List two-factor authentication codes.

Use the `[-i info]` flag to display information about the setup key, it also generates a QR code with the key in URL format that can be scanned by any authenticator.

## Subcommands

- [`sphinx topt add`](https://github.com/5-bare-bones/5bb__sphinx/tree/master/docs/commands/topt/subcommands/add.md): Add a two-factor authentication code.
- [`sphinx topt del`](https://github.com/5-bare-bones/5bb__sphinx/tree/master/docs/commands/topt/subcommands/del.md): Remove two-factor authentication codes from the vault.

## Flags

| Name | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| copy | c | bool | false | Copy code to clipboard |
| info | i | bool | false | Display information about the setup key |
| timeout | t | duration | 0s | Clipboard clearing timeout |

### Timeout units

Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

### Examples

List one and copy to the clipboard:
```
sphinx topt Sample -c
```

List all:
```
sphinx topt
```

Display information about the setup key:
```
sphinx topt Sample -i
```

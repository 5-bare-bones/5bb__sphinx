## Use

`sphinx vault export {manager-name} [-p path]`

## Description

Export entries to other password managers.

This command creates a CSV file with all the entries unencrypted, make sure to delete it after it's used.

Supported password managers:
- 1Password
- Bitwarden
- Keepass/X/XC
- Lastpass

## Flags

|  Name     | Shorthand |     Type      |    Default    |       Description      |
|-----------|-----------|---------------|---------------|------------------------|
| path      | p         | string        | ""            | Destination file path  |

### Examples

Export:
```
sphinx vault export {manager-name} --path {path_to_file}
```

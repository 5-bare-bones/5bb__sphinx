## Use

`sphinx vault backup ([--path {path}] | [--http {url}] [--port {port}])`

## Description

Create vault backup.

## Flags

|  Name     |     Type      |    Default    |                  Description                   |
|-----------|---------------|---------------|------------------------------------------------|
| http      | bool          | false         | Serve the vault file on a http server       |
| path      | string        | ""            | Backup file path                               |
| port      | uint16        | 8080          | Server port                                    |

### Examples

Create file backup:
```
sphinx vault backup --path {path_to_file}
```

Serve vault on a local server:
```
sphinx vault backup --http --port 8080
```

Download vault:
```
curl localhost:8080 > vault_name
sphinx vault restore
```

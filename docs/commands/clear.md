## Use

`sphinx clear  [-c clipboard] [-H history] [-t terminal]`

## Description

Clear clipboard, terminal screen or history.
		
Using the command without passing any flags clears the clipboard and the terminal screen.

## Flags

| Name | Shorthand | Type | Default | Description |
|------|-----------|------|---------|-------------|
| clipboard | c | bool | false | Clear clipboard |
| history | H | bool | false | Remove sphinx commands from terminal history |
| terminal | t | bool | false | Clear terminal screen |

## Examples

Clear terminal and clipboard:
```
sphinx clear
```

Clear clipboard:
```
sphinx clear -c
```

Clear terminal screen:
```
sphinx clear -t
```

Clear sphinx commands from terminal history:
```
sphinx clear -H
```

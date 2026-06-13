## Use

`sphinx file edit <name>  [-e editor] [-l log]`

## Description

Edit a entry.

TODO: change docs to fit entry

Caution: a temporary file is created with a random name, it will be erased right after the first save but it could still be read by a malicious actor.
Notes:
    - Some editors flush the changes to the disk when closed, sphinx won't notice any modifications until then.
    - Modifying the file with a different program will prevent sphinx from erasing the file as its being blocked by another process.

#### Text editors commands
*Editor*: *command*
```
Micro: micro
Vim: vim
Neovim: nvim
Emacs: emacs
Nano: nano
Visual Studio Code: code
Sublime Text: subl
Atom: atom
Coda: coda
Notepad: notepad
Notepad++: notepad++
...
```

#### Image editors commands
*Editor*: *command*
```
GIMP: gimp
Paint: mspaint
Krita: krita
...
```

## Flags

|  Name     | Shorthand |     Type      |    Default    |      Description     |
|-----------|-----------|---------------|---------------|----------------------|
| editor    | e         | string        | ""            | File editor command  |
| log | l | bool | false | Log the temporary file path and wait for modifications |

### Examples

Edit a file:
```
sphinx file edit Sample -e nvim
```

Write a file's content to a temporary file and log its path:
```
sphinx file edit Sample -l
```

# Sphinx - Secure Password Hide IN X? 

[![PkgGoDev](https://pkg.go.dev/badge/github.com/5-bare-bones/5bb__sphinx)](https://pkg.go.dev/github.com/5-bare-bones/5bb__sphinx)

> Created from `kure`
> Modified to be used for learners.

The learners password manager for the command-line that aims to offer a secure and private way of operating with sensitive information by reducing the attack surface to its minimum expression.

## Features

- **Cross-Platform:** Linux, macOS, BSD, Windows and mobile supported.
- **Private:** Self-hosted and completely offline, no connection is established with 3rd parties.
- **Secure:** Each record is encrypted using **AES-GCM** with 256 bit key and a **unique** password derived using Argon2 (**id** version). The user's master password is **never** stored on disk, it's encrypted and temporarily held **in-memory** inside a protected buffer, which is destroyed immediately after use.
- **Sessions:** Run multiple commands by entering the master password only once. They support setting a timeout and running custom scripts.
- **Portable:** Both sphinx and its vault compile to binary files and they can be easily carried around in an external device.
- **Easy-to-use:** Intuitive, does not require advanced technical skills.

## Usage

For further information and examples, visit [docs/commands](/docs/commands).

<img src="https://github.com/user-attachments/assets/64646f5f-a49d-4dea-97d7-99fab2884158" height=600 width=600 />

![Overview](https://user-images.githubusercontent.com/51374959/160211818-b30efbfe-1f1e-44f6-9264-d6faa2f9c0ab.gif)

## Installation

<details><summary>Pre-compiled binaries</summary>
  
Linux, macOS, BSD, Windows and mobile pre-compiled binaries can be downloaded [here](https://github.com/5-bare-bones/5bb__sphinx/releases).

</details>

<details><summary>Scoop (Windows)</summary>

TODO: add scoop bucket to 5-bare-bones

```bash
scoop bucket add 5BB https://github.com/5-bare-bonse/scoop-bucket.git
scoop install 5-bare-bones/5bb__sphinx
```

</details>


<details><summary>Mobile phones terminal emulators</summary>

```bash
curl -LO https://github.com/5-bare-bones/5bb__sphinx/releases/download/{version}/{ARM64 file}
scaphoid archive extract {ARM64 file}
scaphoid file move sphinx $BIN_PATH
```

</details>

<details><summary>Compile from source</summary>

```bash
git clone https://github.com/5-bare-bones/5bb__sphinx
cd sphinx
task build:all 
task install
```

</details>

## Configuration

Out-of-the-box sphinx needs no configuration, it creates a file with the default configuration and the vault at:

- **Linux, BSD**: `$HOME/.sphinx`
- **Darwin**: `$HOME/.sphinx` or `/.sphinx`
- **Windows**: `%USERPROFILE%/.sphinx`

However, to store the configuration file elsewhere or use a different one, set the path to it in the `SPHINX_CONFIG` environment variable.

Head over to the [configuration documentation](/docs/configuration/configuration.md) for a detailed explanation of the configuration file and some [samples](/docs/configuration/samples/).

> [!Note]
> Linux and BSD systems require a utility to write to the clipboard. This could be xsel, xclip, wl-clipboard or the Termux:API add-on.

## Documentation

Learn more about how sphinx works in the [wiki](https://github.com/5-bare-bones/5bb__sphinx/wiki).

## License

This project is licensed under the Apache-2.0 license. See [LICENSE](/LICENSE).

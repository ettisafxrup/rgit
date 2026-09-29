<p align="center">
  <img src="assets/logo.svg" width="128" height="128" alt="rgit logo">
</p>

<h1 align="center">rgit</h1>

<p align="center">
  <strong>git, without the ceremony.</strong><br>
  A fast, friendly command-line companion that turns everyday git chores into single commands.
</p>

<p align="center">
  <img alt="Version" src="https://img.shields.io/badge/version-2.0.0-FFC619?labelColor=153040">
  <img alt="Platforms" src="https://img.shields.io/badge/platform-Windows%20%7C%20Linux%20%7C%20macOS-153040">
  <img alt="Go" src="https://img.shields.io/badge/built%20with-Go-00ADD8">
  <img alt="License" src="https://img.shields.io/badge/license-GPL--3.0-blue">
</p>

---

## Overview

`rgit` wraps the git you already have in a small set of clear, safe commands. Push a change with
one line, see your repository at a glance, undo a commit without memorising `reset` flags, and
generate a `.gitignore`, `LICENSE`, or `README.md` in seconds.

It is a single, self-contained executable with no runtime dependencies. It runs natively in
**Command Prompt, PowerShell, Windows Terminal, Git Bash**, and any Linux or macOS shell.

```text
$ rgit push "Add dark mode"
 • 3 changed files:
   modified   src/theme.css
   modified   src/app.js
   new        src/dark.css
 ? Stage all of these changes? [Y/n]
 ✓ Committed 4f1e02f Add dark mode
 › Pulling latest changes from origin/main
 › Pushing main
 ✓ Pushed main to origin/main
```

## Features

- **One-step publishing.** `rgit push` initialises the repository, connects a remote, stages,
  commits, rebases onto the latest remote work, and pushes, asking before anything important.
- **Safe syncing.** Pulls always use `--rebase --autostash`, so your history stays linear and
  uncommitted work is never lost.
- **A clear overview.** `rgit status` shows the branch, how far it is ahead of or behind the
  remote, the last commit, and every change on one screen.
- **Painless branching.** Create, switch (with an interactive picker), rename (locally *and* on
  the remote), and clean up merged branches.
- **Mistake-friendly.** `rgit undo` removes the last commit and keeps its changes staged. It warns
  you first if that commit was already pushed.
- **Project scaffolding.** Generate a `.gitignore` from 12 presets, a `LICENSE` from 10 popular
  licenses, a complete `README.md`, and an `ARCHITECTURE` tree that follows your `.gitignore`.
- **Helpful everywhere.** It suggests commands when you make typos (`rgit psuh` suggests
  `push`), and `rgit doctor` gives you a fix for every problem it finds.
- **Scriptable.** `--yes` answers every prompt with its default, `NO_COLOR` and `--no-color`
  turn off colours, and exit codes are meaningful.
- **Native Windows installer.** Install per user (no admin rights) or for all users. The
  installer puts rgit on your `PATH`, and the uninstaller removes it cleanly.

## Installation

Every release on the [Releases page](https://github.com/ettisafxrup/rgit/releases/latest)
contains these files:

| File | For |
| --- | --- |
| `rgit-windows-x64-setup.exe` | Windows 10/11: the installer (recommended) |
| `rgit-windows-x64.exe` | Windows: a portable executable that needs no installation |
| `rgit-linux-x64` | Linux on Intel/AMD (x86_64) |
| `rgit-linux-arm64` | Linux on ARM64 (for example Raspberry Pi 4/5 and AWS Graviton) |
| `rgit-macos-arm64` | macOS on Apple Silicon (M1 and later) |
| `rgit-macos-x64` | macOS on Intel |
| `SHA256SUMS.txt` | Checksums for verifying the downloads |

Website: [ettisafxrup.github.io/rgit](https://ettisafxrup.github.io/rgit/)

### Windows: installer (recommended)

1. Download [`rgit-windows-x64-setup.exe`](https://github.com/ettisafxrup/rgit/releases/latest/download/rgit-windows-x64-setup.exe).
2. Run it and choose **Install for me only** (no administrator rights needed) or
   **Install for all users**.
3. Keep **Add rgit to the PATH** ticked.
4. Open a **new** terminal and check that rgit works:

   ```powershell
   rgit version
   rgit doctor
   ```

If an older, shell-based rgit (1.x) is installed, an all-users installation removes it
automatically. A per-user installation tells you to remove it from **Settings → Apps**.

**Silent / unattended installation** (for IT or scripts):

```powershell
rgit-windows-x64-setup.exe /VERYSILENT /SUPPRESSMSGBOXES /CURRENTUSER   # per user
rgit-windows-x64-setup.exe /VERYSILENT /SUPPRESSMSGBOXES /ALLUSERS      # all users (elevated)
rgit-windows-x64-setup.exe /VERYSILENT /TASKS=""                        # without touching PATH
```

To uninstall, use **Settings → Apps → rgit**, or run `unins000.exe` from the installation
folder. Both remove rgit from your `PATH` as well.

**Portable:** rename `rgit-windows-x64.exe` to `rgit.exe` and put it in any folder on your `PATH`.

### Linux and macOS

Install the latest release with one command. The script detects your system and processor,
verifies the checksum, and installs to `/usr/local/bin`, or to `~/.local/bin` if that folder
is not writable:

```bash
curl -fsSL https://raw.githubusercontent.com/ettisafxrup/rgit/main/scripts/install.sh | sh
```

To install manually, download the file for your platform and put it on your `PATH`:

```bash
curl -fLo rgit https://github.com/ettisafxrup/rgit/releases/latest/download/rgit-linux-x64
chmod +x rgit
sudo mv rgit /usr/local/bin/
rgit version
```

On macOS, a manually downloaded binary is quarantined because it is not notarized. Allow it once
with `xattr -d com.apple.quarantine /usr/local/bin/rgit`. The install script does this for you.

### With Go

If you have Go 1.24 or newer:

```bash
go install github.com/ettisafxrup/rgit/cmd/rgit@latest
```

### Requirements

- [git](https://git-scm.com) 2.23 or newer, available on your `PATH`
- Windows 10/11 (x64 or ARM64), Linux, or macOS

## Getting started

```bash
rgit help              # all commands, grouped
rgit help push         # details and examples for one command
rgit login             # tell git who you are (name and email)
rgit doctor            # check your setup
```

A typical day:

```bash
rgit pull                        # start from the latest code
rgit new feature/search          # branch off
# ... work ...
rgit status                      # what have I changed?
rgit push "Add search box"       # stage + commit + rebase + push
rgit switch main                 # back to main
rgit cleanup                     # delete branches that are merged
```

Starting a brand-new project:

```bash
mkdir my-app && cd my-app
rgit ignore node vscode windows  # .gitignore from presets
rgit license mit                 # LICENSE with your name and the current year
rgit readme                      # README.md with structure, setup and license sections
rgit push "Initial commit"       # creates the repo, asks for the remote URL, pushes
```

## Command reference

### Everyday

| Command | What it does |
| --- | --- |
| `rgit push [message]` | Stages, commits, rebases onto the remote, and pushes. Quotes around the message are optional. It creates the repository and remote when they are missing. Alias: `p`. |
| `rgit pull [branch]` | Pulls with rebase and autostash. In an empty folder, it offers to connect the folder to a remote and download it. |
| `rgit sync` | Pulls remote changes, then pushes your commits, without creating a new commit. |
| `rgit status` | Shows the branch, its upstream, ahead and behind counts, the last commit, changes, and stashes. Alias: `st`. |
| `rgit log [-n N] [--all]` | Shows a compact, colourful commit graph. Alias: `l`. |
| `rgit undo` | Undoes the last commit and keeps its changes staged. It warns you if the commit was already pushed. |

### Branches

| Command | What it does |
| --- | --- |
| `rgit new <branch>` | Creates a branch and switches to it. Alias: `branch`. |
| `rgit switch [branch]` | Switches branches. Without a name, it shows a numbered picker. Aliases: `sw`, `checkout`. |
| `rgit rename <new>` <br> `rgit rename <old>:<new>` <br> `rgit rename <old> <new>` | Renames a branch. If the branch exists on the remote, rgit pushes the new name first and only then deletes the old one. Alias: `mv`. |
| `rgit cleanup [base] [--prune]` | Deletes local branches that are fully merged into the base branch. `main`, `master`, `develop`, and your current branch are always kept. Alias: `clean`. |

### Project files

| Command | What it does |
| --- | --- |
| `rgit ignore [preset \| pattern]...` | Adds presets (`go`, `node`, `python`, `java`, `rust`, `dotnet`, `cpp`, `windows`, `macos`, `linux`, `vscode`, `jetbrains`) or patterns to `.gitignore`. It never adds duplicates. Without arguments, it opens a file picker. Use `--list` to see the presets. |
| `rgit license [id] [--name N] [--year Y]` | Writes a `LICENSE`. Available IDs: `mit`, `apache-2.0`, `gpl-3.0`, `bsd-3-clause`, `bsl-1.0`, `agpl-3.0`, `lgpl-3.0`, `mpl-2.0`, `unlicense`, `cc-by-4.0`. |
| `rgit readme` | Generates a `README.md`. It reads your GitHub details from the remote and offers to add a license if the project has none. |
| `rgit archi [-o file] [--depth N] [--print]` | Writes the project tree to `ARCHITECTURE`. Inside a repository, the tree follows `.gitignore`. Alias: `tree`. |

### Setup

| Command | What it does |
| --- | --- |
| `rgit clone <repo> [folder]` | Clones a repository. Accepts the GitHub shorthand `owner/repo`. |
| `rgit login` | Sets your global git name and email (validated), then checks access to the remote. |
| `rgit whoami` | Shows the identity used for your commits. Aliases: `userinfo`, `me`. |
| `rgit open [remote]` | Opens the repository page in your browser. On GitHub, it opens the current branch. Alias: `web`. |
| `rgit doctor` | Checks git, your identity, the credential helper, `PATH`, the branch, and the remote. |
| `rgit version` | Prints the version, platform, and Go runtime. |
| `rgit help [command]` | Shows the command overview, or detailed help for one command. |

### Global flags

| Flag | Meaning |
| --- | --- |
| `-y`, `--yes` | Answer yes to every confirmation and accept default values (useful in scripts) |
| `--no-color` | Disable coloured output. The `NO_COLOR` environment variable is also honoured. |
| `-h`, `--help` | Show help for the command |
| `-v`, `--version` | Print the version |

Flags can go anywhere: `rgit push -y "msg"` and `rgit push "msg" -y` both work. Everything
after `--` is treated as text, for example `rgit push -- "-1 bug"`.

**Exit codes:** `0` success · `1` error or aborted · `2` invalid usage.

## What's new in 2.0

rgit 2.0 is a complete rewrite in Go. Version 1.x was a set of Bash scripts, so it could not run
in Command Prompt or PowerShell. Version 2.0 is a native executable for every platform.

It also fixes these bugs from 1.x:

- `push` no longer asks you to type a different branch name after committing. It always pushes
  the branch you are on.
- The first push to a remote that already has commits (a README created on GitHub, for example)
  now works. rgit matches the remote's default branch (`main`, `master`, …) and rebases onto it.
- An empty commit no longer crashes the flow. rgit reports "Nothing new to commit" and continues
  with the push.
- Repositories are detected correctly from any subfolder, not only from the root.
- `rename` no longer assumes the remote is called `origin`, and it never deletes the remote branch
  before the new one has been pushed.
- `cleanup` no longer assumes the base branch is `main`. It also no longer deletes branches whose
  names merely *contain* "main", and it no longer fails when there is nothing to clean.
- `ignore` no longer adds `.gitignore` to itself and never writes duplicate entries.
- `license` works with names that contain `/` or `&`, and it asks before overwriting.
- `readme` asks before overwriting and no longer creates an `ARCHITECTURE` file as a side
  effect. The typos in the generated text are fixed too.
- Error messages go to stderr with the correct colours, and every failure returns a non-zero exit
  code.

## Project structure

```text
.
├── cmd/rgit/             entry point (main.go) and embedded Windows resources
├── internal/
│   ├── cli/              command model, argument parsing, help and suggestions
│   ├── commands/         one file per command: push.go, pull.go, rename.go, …
│   ├── git/              thin, tested wrapper around the git executable
│   ├── templates/        embedded licenses, .gitignore presets, README template
│   ├── tree/             project tree rendering
│   ├── ui/               colours, Windows console setup, prompts
│   ├── testenv/          helpers for tests that run real git
│   └── version/          build version
├── installer/            Inno Setup script and wizard artwork
├── scripts/              build.ps1, build.sh, install.sh, Windows resources
├── tools/genlogo/        generates the logo, icon, installer and website images
├── docs/                 the project website (index.html), ready for GitHub Pages
├── assets/               logo.svg, logo.png, rgit.ico
└── shell-rgit/           the original 1.x shell implementation, kept for reference
```

rgit calls the real `git` executable instead of reimplementing it. Your credential helpers,
hooks, aliases, and configuration keep working exactly as they do with git.

## Building from source

Requirements: Go 1.24+, git. On Windows, you also need [Inno Setup](https://jrsoftware.org/isinfo.php)
6 or newer to build the installer. MinGW `windres` is optional and is only needed to rebuild the
embedded icon.

```powershell
# Windows: runs vet and tests, then builds the Windows exe and installer into dist\
.\scripts\build.ps1

# A complete release: every platform plus SHA256SUMS.txt
.\scripts\build.ps1 -All

# Regenerate the logo, icon and installer artwork
.\scripts\build.ps1 -Assets
```

```bash
# Linux / macOS
./scripts/build.sh          # build for this machine
./scripts/build.sh --all    # build for every platform
```

To run the test suite on its own:

```bash
go test ./...
```

The tests create temporary repositories and local "remotes", and they use an isolated git
configuration, so your own settings are never read or changed.

## Contributing

Contributions are welcome, whether a bug fix, a new command, or a documentation improvement.

1. Fork the repository and create a branch: `rgit new feature/my-change`
2. Make your change and add tests next to the code you touched.
3. Make sure `go vet ./...` and `go test ./...` pass.
4. Push the branch (`rgit push "Describe the change"`) and open a pull request.

Please keep each pull request focused on a single change. To report a bug or suggest a feature,
[open an issue](https://github.com/ettisafxrup/rgit/issues/new).

## License

rgit is licensed under the **GNU General Public License v3.0**. See [LICENSE](./LICENSE) for the full text.

---

<p align="center">
  Made by <a href="https://github.com/ettisafxrup">@ettisafxrup</a>
</p>

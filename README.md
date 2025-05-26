# rgit — Git Command Macro Tool

**rgit** is a lightweight command-line macro tool designed to speed up and simplify your Git workflow. It provides shortcuts for common Git commands like push, pull, and branch rename — saving you time and keystrokes.

---

## Features

- Simplified Git command aliases
- Color-coded terminal output for better readability
- Interactive prompts for safer operations
- Easily extendable with custom scripts

---

## Installation

Download the installer from <a href="https://github.com/ettisafxrup/rgit/blob/main/installer/rgit.exe" download>here</a> and install.
<br>
That's it, you're just ready to go! Enjoy <b>rgit</b> commands directly in your terminal!

---

## Usage

Open your terminal (Git Bash, Linux shell, or Windows CMD with the proper environment) and run:

```bash
rgit help
```

## Commands

- push <br>

  It automatically checks for initialization, initializes if it's not, creates branches and commits and directly pushes to the repository in just a single prompt!

```bash
rgit push
```

- pull <br>

  It automatically checks for initialization and pulls the repository to the local work setup - likely a ready-to-work command.

```bash
rgit pull
```

- rename <br>

  Changing a branch's name is way more easy with rgit. If I want to change a branch name from "test" to "dev, all I need to do is.

```bash
# rgit rename <current_branch_name>:<new_branch_name>
# as an example:

rgit rename test:dev
```

- ignore <br>

  create an .gitignore file by using rgit, more easily and way more fast.

```bash

rgit ignore

# It'll be asking file names, add one by one manually and hit enter when you're done.
```

It will give you all the list it can perform in the terminal.

# Troubleshooting

- Make sure Git is installed and available in your terminal. (Having git installed is a must to use rgit)

- Restart terminal after installation to refresh environment variables.

- If Windows Prompt having issues while running it, Use Git Bash on Windows for best compatibility.

###

# License

MIT License. See <a href="">LICENSE</a> file for details.

###

# Contributing

<b>rgit</b> is an open sourced toolkit. Feel free to contribute, open issues or submit pull requests!

###

# Author

Ettisaf Rup
<br>
Author, rgit toolkit.
<br>
Ug, KUET CSE.

# rgit — Git Command Macro Tool

`rgit` is a lightweight command-line macro tool designed to speed up and simplify your Git workflow. It provides shortcuts for common Git commands like push, pull, and branch rename — saving you time and keystrokes.

---

📦 **Current Version: v1.1.0**

- 🤷‍♀️ Patch update: Minor bug fixes
- ✨ New Features: `rgit readme`, `rgit archi`, `rgit license`
- 🔄 Improved: updated `gitignore` generation (`rgit ignore`)

---

## 🦋 Features

- 🕶 Simplified Git command aliases
- 🎒 Ready-made resources for project initialization.
- 🎨 Color‑coded terminal output for better readability
- 🔒 Interactive prompts for safer operations
- 🔧 Easily extendable with custom scripts

---

## 📥 Installation

Download the installer from <a href="https://github.com/ettisafxrup/rgit/blob/main/installer/rgit.exe" download>here</a> and install.  
That's it, you're just ready to go! Enjoy <b>rgit</b> commands directly in your terminal!

---

## 🚀 Kickstart

Open your terminal (Git Bash, Linux shell, or Windows CMD with the proper environment) and run:

```bash
rgit help
```

It'll be showing you all the updated commands in the terminal.

## ⚙ Commands

✅ Here's a list of available rgit commands:

```text
 > rgit push
   — Automates init → commit → push
   # rgit push "your commit message"

 > rgit pull
   — Pulls latest remote changes and updates local branch
   # rgit pull

 > rgit rename
   — Renames current branch to a new one
   # rgit rename current:new

 > rgit cleanup
   — Renames current branch to a new one
   # rgit rename current:new

 > rgit ignore
   — Creates a smart .gitignore file
   # rgit ignore

 > rgit readme
   — Generates a basic README.md file
   # rgit readme

 > rgit license
   — Adds a LICENSE file (MIT by default)
   # rgit license

 > rgit archi
   — Generates a project structure in a file named ARCHITECTURE
   # rgit archi

```

## 📁 Project Structure — rgit

```text

📁 Application Structure:

.
├── 📂 assets
│   ├── favicon.ico
│   ├── working_rgit.iss
├── 📂 configs
│   ├── color.sh
│   ├── msg.sh
├── 📂 installer
│   ├── rgit.exe
├── LICENSE
├── 📂 prompts
│   ├── 📂 core
│   │   ├── cleanup.sh
│   │   ├── pull.sh
│   │   ├── push.sh
│   │   ├── rename.sh
│   ├── 📂 create
│   │   ├── archi.sh
│   │   ├── forge.sh
│   │   ├── ignore.sh
│   │   ├── 📂 LICENSE
│   │   │   ├── apache-2.0.txt
│   │   │   ├── bsd-3-clause.txt
│   │   │   ├── gpl-3.0.txt
│   │   │   ├── mit.txt
│   │   │   ├── unlicense.txt
│   │   ├── license.sh
│   │   ├── readme.sh
│   ├── help.sh
│   ├── 📂 info
│   │   ├── login.sh
│   │   ├── userinfo.sh
│   ├── pull.sh
│   ├── push.sh
│   ├── rename.sh
├── rgit (main_script)

```

## 🌸 Cloning Repository

You can _fork_ the repo, or you can directly clone this repository locally.

🔧 Run this in your terminal:

```bash
git clone https://github.com/ettisafxrup/rgit
```

## 🤝 Contributing

Feel free to contribute, we love contributions from the community! ❤  
Whether it's fixing a bug or adding a feature, or improving this documentation — every it of help counts.

🤔 **How to Contribute**

---

1. **Fork the repository**  
   _forking the repository will be working, cloning locally won't work here._

2. **Create a new branch**

```bash
git checkout -b feature/your-feature-name
```

3. **Make and Commit your changes**

```bash
git commit -m "your Commit Message"
```

4. **Push your Fork**

```bash
git push origin feature/your-feature-name
```

5. **Open a pull request**

We'll review your PR and get back to you soon!

---

## 📚 Guidelines

- Make sure your code follows the project's coding style.
- Keep PRs focused and minimal. One change per pull request will be the best.
- Be respectful and constructive in taking major decisions.

## ⚒ Issues

Found a bug? Have a feature request?  
Please <a href="https://github.com/ettisafxrup/rgit/issues/new" target="_blank">open an issue</a> — we'd love to hear from you.

## 📄 License

The project is under the license of MIT.

> You are free to use, copy, odify and distribute this source.  
> See the [LICENSE](./LICENSE) file for full legal terms.

---

> **GitHub**: [@ettisafxrup](https://github.com/ettisafxrup)  
> **Repository**: [https://github.com/ettisafxrup/rgit](https://github.com/ettisafxrup/rgit)

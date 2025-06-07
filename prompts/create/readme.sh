#!/bin/bash

# IMPORTS:
# Terminal Beauty
source "$RCOMMANDS_ROOT/configs/color.sh"
source "$RCOMMANDS_ROOT/configs/msg.sh"

source "$RCOMMANDS_ROOT/prompts/create/forge.sh" # Project Structure
source "$RCOMMANDS_ROOT/prompts/create/license.sh" # Project License

readme(){
    echo -e "${BOLD}${BLUE}README generator started!${NC}"
    echo -e "${INFO} Generating Project Structure... \n"
    # Creating a forge of this project
    forge
    # Creating License if there is no license exists.
    if [ -f "LICENSE" ]; then
        LICENSE_NAME=$(head -n 1 LICENSE)
        echo -e "${INFO} ${BOLD}${RED}LICENSE${NC} file found:${YELLOW} ${LICENSE_NAME} ${NC}"
    else
        license  # Call your license function or prompt here
    fi

echo -e "\n${BOLD}Taking Project Details...${NC}"

read -ep "${GREEN}${BOLD}>${NC} Project Title: ${NC}" TITLE
read -ep "${GREEN}${BOLD}>${NC} Description: ${NC}" DESCRIPTION
read -ep "${GREEN}${BOLD}>${NC} Project Run Command (eg. npm run dev): ${NC}" RUN
echo -e "\n\t ${TICK} Project Details Initialized!\n"
echo -e "${GREEN}#${NC} ${BOLD}Github Account:${NC}"
read -ep "  ${SKY_BLUE}GitHub Username: ${NC}" USERNAME
read -ep "  ${YELLOW}Repository Name: ${NC}" REPO

cat > README.md <<EOF
# $TITLE 🚀

$DESCRIPTION
<!-- Description, $TITLE -->
---

## 🏃‍♂️ Run  
Open up your terminal and write this command:  

\`\`\`bash
$RUN
\`\`\`

## 💖 Features

- F.1
- F.2
- F.3
(optional, you can keep it if you want.)

## 📁 Project Structure
\`\`\`text

$(cat ARCHITECTURE)

\`\`\`

## 🦋 Cloning Repository

You can _fork_ the repo, or you can directly clone this repository locally.  

🔧 Run this in your terminal:  

\`\`\`bash
git clone https://github.com/$USERNAME/$REPO
\`\`\`

## 🤝 Contributing

Feel free to contribute, we love contributions from the community! ❤  
Whether it's fixing a bug or adding a feature, or improving this documentation - every it of help counts.  

🤔 **How to Contribute**
---  

1. **Fork the repository**  
*forking the repository will be working, cloning locally won't work here.*  

2. **Create a new branch**  

\`\`\`bash
git checkout -b feature/your-feature-name
\`\`\`
3. **Make and Commit your changes**  

\`\`\`bash
git commit -m "Myour Commit Message"
\`\`\`
4. **Push your Fork**  
\`\`\`bash
git push origin feature/your-feature-name
\`\`\`
5. **Open a pull request**  
 
We'll review your PR and get back to you soon!  


---

## 📚 Guidelines
- Make sure ytour code follows the project's coding style.
- Keep PRs focused and minimal. One change per pull request will be the best.
- Be respectful and constructive in taking decisions.  

## ⚒ Issues
Found a bug? Have a feautre request?  
Please <a href="https://github.com/$USERNAME/$REPO/issues/new" target="_blank">open an issue</a> — we'd love to hear from you.

## 📄 License  

**[LICENSE](./LICENSE) ATTACHED**  
> You are free to use, copy, odify and distribute this source.  
> See the [LICENSE](./LICENSE) file for full legal terms. 

---

> **GitHub**: [@$USERNAME](https://github.com/$USERNAME)  
> **Repository**: [https://github.com/$USERNAME/$REPO](https://github.com/$USERNAME/$REPO)

EOF

echo -e "\n${BLUE}${BOLD}README.md has been created in the current directory.${NC}"
}
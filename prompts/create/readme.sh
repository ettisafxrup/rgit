#!/bin/bash

# IMPORTS:
# Terminal Beauty
source "$RCOMMANDS_ROOT/configs/color.sh"
source "$RCOMMANDS_ROOT/configs/msg.sh"


readme(){

echo "Let's generate a README.md for your GitHub project!"

read -ep "${INFO}${BOLD} Project Title: ${NC}" TITLE
read -ep "${INFO}${BOLD} Description: ${NC}" DESCRIPTION
read -ep "${INFO}${BOLD} Source: ${NC}" SOURCE
read -ep "${TICK}${BOLD} Usage Information: ${NC}" USAGE
read -ep "${TICK}${BOLD} Contribution Guidelines: ${NC}" CONTRIB
read -ep "${TICK}${BOLD}${YELLOW} License (e.g., MIT, Apache-2.0, GPL-3.0): ${NC}" LICENSE
read -ep "${SKY_BLUE} GitHub Username: ${NC}" USERNAME
read -ep "${YELLOW} Repository Name: ${NC}" REPO

cat > README.md <<EOF
# $TITLE 🚀

$DESCRIPTION

---

## 📦 Source

$SOURCE

## 🛠️ Usage

\`\`\`bash
$USAGE
\`\`\`

## 🤝 Contributing

$CONTRIB

## 📄 License

This project is licensed under the **$LICENSE** license.

---

> GitHub: [@$USERNAME](https://github.com/$USERNAME)  
> Repo: [https://github.com/$USERNAME/$REPO](https://github.com/$USERNAME/$REPO)

EOF

echo -e "\n ${BLUE}${BOLD}README.md has been created in the current directory.${NC}"
}
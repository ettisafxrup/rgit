#!/bin/bash
# ==== Terminal Message Templates ====
source "$RCOMMANDS_ROOT/configs/msg.sh"

pull(){

echo -e "${INFO} Checking if this is a Git repository..."

if [ ! -d ".git" ]; then
    echo -e "${RED}${BOLD} This is not a Git repository. ${NC}"
    git init
fi

# Detect remote and current branch
REMOTE_NAME=$(git remote | head -n 1)
CURRENT_BRANCH=$(git symbolic-ref --short HEAD)

if [ -z "$REMOTE_NAME" ]; then
    echo -e "$NO_GIT_REMOTE_FOUND"
    read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
    read -rp "$PROMPT_REPO_URL" REPO_URL

        # FILTERING REPO_URL:
    # Remove trailing question marks
    REPO_URL="${REPO_URL%%\?}"
    # Remove trailing spaces (and tabs)
    REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"

    git remote add "$REMOTE_NAME" "$REPO_URL"
fi

echo -e "${INFO} Pulling latest changes from ${YELLOW}$REMOTE_NAME/$CURRENT_BRANCH${NC}..."

# Perform git pull with rebase
if git pull --rebase "$REMOTE_NAME" "$CURRENT_BRANCH"; then
    echo -e "${SUCCESS} Rebase and pull successful."
else
    echo -e "${ERROR} Git pull failed. Please resolve any conflicts and run:"
    echo -e "${YELLOW}    git rebase --continue${NC}  or  ${YELLOW}git rebase --abort${NC}"
    echo "$FAIL"
    exit 1
fi


}
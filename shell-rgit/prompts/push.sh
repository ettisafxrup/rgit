#!/bin/bash

# SAFETY FLAGS
set -euo pipefail
IFS=$'\n\t'

# IMPORTS:
source "$RCOMMANDS_ROOT/configs/msg.sh"
COMMIT_MSG="$2"
echo "$COMMIT_MSG"

push() {

    if [ -d ".git" ]; then
        echo -e "${BLUE}${BOLD}[rgit]: directory found initialized!${NC}"

        # Check for remote
        if [ -z "$(git remote)" ]; then
            echo -e "$NO_GIT_REMOTE_FOUND"
            read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
            read -rp "$PROMPT_REPO_URL" REPO_URL

            # Clean URL
            REPO_URL="${REPO_URL%%\?}"
            REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"

            git remote add "$REMOTE_NAME" "$REPO_URL"
        else
            REMOTE_NAME=$(git remote | head -n 1)
        fi

        # Ask before staging everything
        read -rp "Do you want to stage all changes? [y/N]: " CONFIRM
        if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
            echo -e "${RED}Aborted staging. Exiting.${NC}"
            exit 1
        fi
        git add .

        # Commit
        if [ -z "$COMMIT_MSG" ]; then
            echo -e "$COMMIT_EMPTY"
            git commit -m "pushed on $TODAY"
        else
            git commit -m "$COMMIT_MSG"
        fi

        CURRENT_BRANCH=$(git symbolic-ref --short HEAD)

        echo -e "${INFO} Pulling latest changes with rebase..."
        git pull --rebase "$REMOTE_NAME" "$CURRENT_BRANCH"

        branch_count=$(git branch | wc -l)
        if [ "$branch_count" -gt 1 ]; then
            echo "${INFO} Multiple branches found:"
            git branch
            read -rp "${BOLD}Write your branch to push: " CURRENT_BRANCH
            echo "${TICK}"
        fi

        # Check for upstream
        if git rev-parse --abbrev-ref --symbolic-full-name @{u} &>/dev/null; then
            git push
        else
            git push -u "$REMOTE_NAME" "$CURRENT_BRANCH"
        fi

    else
        echo -e "${YELLOW}No Git repo found. Initializing...${NC}"
        git init
        read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
        read -rp "$PROMPT_REPO_URL" REPO_URL

        # Clean URL
        REPO_URL="${REPO_URL%%\?}"
        REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"

        git remote add "$REMOTE_NAME" "$REPO_URL"

        CURRENT_BRANCH="main"

        echo -e "${INFO} Pulling remote into new repo..."
        if ! git pull --rebase --allow-unrelated-histories "$REMOTE_NAME" "$CURRENT_BRANCH"; then
            echo -e "${RED}Failed to pull from remote. Resolve conflicts manually.${NC}"
            exit 1
        fi

        read -rp "Do you want to stage all changes? [y/N]: " CONFIRM
        if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
            echo -e "${RED}Aborted staging. Exiting.${NC}"
            exit 1
        fi
        git add .

        if [ -z "$COMMIT_MSG" ]; then
            echo -e "$COMMIT_EMPTY"
            git commit -m "Push: $TODAY"
        else
            git commit -m "$COMMIT_MSG"
        fi

        git push -u "$REMOTE_NAME" "$CURRENT_BRANCH"
        echo -e "$SUCCESS"
    fi
}

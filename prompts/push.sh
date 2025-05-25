#!/bin/bash

# IMPORTS:
# Terminal Messages
source "$RCOMMANDS_ROOT/configs/msg.sh"

push() {
    COMMIT_MSG="$1"

    if [ -d ".git" ]; then
        echo -e "${BLUE}${BOLD}[rgit]: directory found initialized!${NC}"

        # Check if remote exists; if not, ask for one and add it
        if [ -z "$(git remote)" ]; then
            echo -e "$NO_GIT_REMOTE_FOUND"
            read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
            read -rp "$PROMPT_REPO_URL" REPO_URL

            # FILTERING REPO_URL
            REPO_URL="${REPO_URL%%\?}"                                # Remove trailing question marks
            REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"       # Remove trailing spaces (and tabs)

            git remote add "$REMOTE_NAME" "$REPO_URL"
        else
            REMOTE_NAME=$(git remote | head -n 1)  # Use first remote found
        fi

        git add .

        if [ -z "$COMMIT_MSG" ]; then
            echo -e "$COMMIT_EMPTY"
            git commit -m "pushed on $TODAY"
        else
            git commit -m "$COMMIT_MSG"
        fi

        # Get the current branch name
        CURRENT_BRANCH=$(git symbolic-ref --short HEAD)

        # Pull remote changes with rebase to avoid conflicts
        git pull --rebase "$REMOTE_NAME" "$CURRENT_BRANCH"

        # Check if multiple branches exists

        if [ "$branch_count" -gt 1 ]; then
            echo "${INFO} Multiple branches found:"
            git branch
            read -rp "${BOLD}Write your branch:" CURRENT_BRANCH
            echo "${TICK}"

        fi

        # Check if upstream is set
        UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
        if [ -z "$UPSTREAM" ]; then
            git push -u "$REMOTE_NAME" "$CURRENT_BRANCH"
        else
            git push
        fi

    else
        # This code will be executed if Git is not initialized in the directory
        git init
        read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
        read -rp "$PROMPT_REPO_URL" REPO_URL

        # FILTERING REPO_URL
        REPO_URL="${REPO_URL%%\?}"                                # Remove trailing question marks
        REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"       # Remove trailing spaces (and tabs)

        git remote add "$REMOTE_NAME" "$REPO_URL"

        CURRENT_BRANCH="main"  # Default to main when initializing
        git pull --rebase "$REMOTE_NAME" "$CURRENT_BRANCH"

        git add .

        if [ -z "$COMMIT_MSG" ]; then
            echo -e "$COMMIT_EMPTY"
            git commit -m "Push: $TODAY"
        else
            git commit -m "$COMMIT_MSG"
        fi

        git push -u "$REMOTE_NAME" main
        echo -e "$SUCCESS"
    fi
}

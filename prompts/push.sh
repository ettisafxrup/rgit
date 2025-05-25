#!/bin/bash

# IMPORTS:
# Terminal Messages
source "$RCOMMANDS_ROOT/configs/msg.sh"


push(){
    COMMIT_MSG="$1"
if [ -d ".git" ]; then
        echo -e "$GIT_INIT_FOUND"

        # Check if remote exists; if not, ask for one and add it
        if [ -z "$(git remote)" ]; then
            echo -e "$NO_GIT_REMOTE_FOUND"
            read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
            read -rp "$PROMPT_REPO_URL" REPO_URL

                # FILTERING REPO_URL:
            # Remove trailing question marks
            REPO_URL="${REPO_URL%%\?}"
            # Remove trailing spaces (and tabs)
            REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"

            git remote add "$REMOTE_NAME" "$REPO_URL"
        else
            REMOTE_NAME=$(git remote | head -n 1)  # Use first remote found
        fi

        # Add changes and commit if message provided

        git add .

        if [ -z "$COMMIT_MSG" ]; then
            echo -e "$COMMIT_EMPTY"
            git commit -m "pushed on $TODAY"
        else
            git commit -m "$COMMIT_MSG"
        fi

        # Getting the current branch name
        CURRENT_BRANCH=$(git symbolic-ref --short HEAD)

        # Pull remote changes with rebase to avoid conflicts
        git pull --rebase "$REMOTE_NAME" "$CURRENT_BRANCH"

        # Push and set upstream if it’s first push
        # Checking if there is any, if yes then going for a straight git push (at else block).

        UPSTREAM=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null)
        echo git branch -M main
        if [ -z "$UPSTREAM" ]; then
            git push -u "$REMOTE_NAME" "$CURRENT_BRANCH"
        else
            git push
        fi

    else
        # This code will be executed if Git is not initialized at the directory (fresh rgit run)
        git init
        read -rp "$PROMPT_REMOTE_NAME" REMOTE_NAME
        read -rp "$PROMPT_REPO_URL" REPO_URL

            # FILTERING REPO_URL:
        # Remove trailing question marks
        REPO_URL="${REPO_URL%%\?}"
        # Remove trailing spaces (and tabs)
        REPO_URL="${REPO_URL%"${REPO_URL##*[![:space:]]}"}"
            
        git remote add "$REMOTE_NAME" "$REPO_URL"
        git pull --rebase "$REMOTE_NAME" "$CURRENT_BRANCH"

        git add .
        if [ -z "$COMMIT_MSG" ]; then
            echo -e "$COMMIT_EMPTY"
            git commit -m "Push: $TODAY"
        else
            git commit -m "$COMMIT_MSG"
        fi
        git branch -M main
        git push -u "$REMOTE_NAME" main
    fi

    echo -e "$SUCCESS"
}
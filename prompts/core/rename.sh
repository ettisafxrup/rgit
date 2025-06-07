#!/bin/bash
source "$RCOMMANDS_ROOT/configs/color.sh"
source "$RCOMMANDS_ROOT/configs/msg.sh"

rename(){
    INPUT="$1"

    # Set delimiter as '->' and read into two variables
    IFS=':' read -r CURRENT_BRANCH_NAME NEW_BRANCH_NAME <<< "$INPUT"

    if [[ -z "$CURRENT_BRANCH_NAME" || -z "$NEW_BRANCH_NAME" ]]; then
    echo -e "${ERROR} ${RED} Invalid input.${NC} ${BOLD} Use format: current_branch:new_branch${NC}"
    exit 1
    fi

    echo "${BLUE}${BOLD} Current Branch: $CURRENT_BRANCH_NAME"
    echo "${INFO}${BOLD} Renaming $CURRENT_BRANCH_NAME to $NEW_BRANCH_NAME"
    read -rp "${TICK} ${BOLD}Press Enter to proceed...${NC}"

    git branch -m "$CURRENT_BRANCH_NAME" "$NEW_BRANCH_NAME"
    git push origin --delete "$CURRENT_BRANCH_NAME"
    git push origin "$NEW_BRANCH_NAME"
    git push --set-upstream origin "$NEW_BRANCH_NAME"

    echo "${TICK} Branch named to $NEW_BRANCH_NAME Successfully."
}
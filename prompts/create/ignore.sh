#!/bin/bash

# IMPORTS:
# Terminal Beauty
source "$RCOMMANDS_ROOT/configs/color.sh"
source "$RCOMMANDS_ROOT/configs/msg.sh"

ignore() {
    GITIGNORE_FILE=".gitignore"
    FILE_COUNT=0

    echo "${BLUE}${BOLD}[help]:${NC} no different of adding filenames in git (eg. file.txt, *.exe) ${NC}"
    echo -e "Enter ${BOLD}${BLUE}FILENAMES${NC} to ${BOLD}${RED}IGNORE${NC}.\nPress blank ${BOLD}${RED}ENTER${NC} to finish:${NC}"

    while true; do
        read -rp "${BOLD}${GREEN}  > ${NC}" filename
        if [ -z "$filename" ]; then
            break
        fi
        echo "$filename" >> "$GITIGNORE_FILE"

        ((FILE_COUNT++))
    done

    echo -e "\n $FILE_COUNT ${BOLD}FILENAMES${NC} added successfully."
    echo -e "${TICK} ${RED}.gitignore${NC} Initialized in this directory."
}

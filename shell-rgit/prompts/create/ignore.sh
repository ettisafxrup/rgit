#!/bin/bash
# IMPORTS:
# Importing Help Command
source "$RCOMMANDS_ROOT/configs/msg.sh" # as msg
source "$RCOMMANDS_ROOT/prompts/help.sh" # as help

ignore() {
    GITIGNORE_FILE=".gitignore"
    FILE_COUNT=0
    DIR_COUNT=0

    if [ ! -f ".gitignore" ]; then
        echo ".gitignore" >> "$GITIGNORE_FILE"
    fi

    echo -e "${BLUE}${BOLD}[help]:${NC} no different of adding filenames in git (eg. file.txt, *.exe) ${NC}"
    echo -e "${YELLOW}${BOLD}[HIGHLIGHT]:${NC} To Ignore a Directory, Add \"/\" before Folder Name. (eg. /assets, /env ${NC}"

    # OLD CODE: Manually Adding Files
    # while true; do
        
    #     read -rp "${BOLD}${GREEN}  > ${NC}" filename
    #     if [ -z "$filename" ]; then
    #         break
    #     fi
    #     echo "$filename" >> "$GITIGNORE_FILE"
    #     ((FILE_COUNT=FILE_COUNT+1))
    # done

    echo -e "\n Select ${BOLD}${BLUE}FILENAMES${NC} to ${BOLD}${RED}IGNORE${NC}.\n\n${BOLD}${RED}0) EXIT${NC}${NC}"
    select filename in *; do
        while true; do
            if [ -z "$filename" ]; then
                break 2  # Exit both select and while
            fi
        # Check if it's a directory
            if [ -d "$filename" ]; then
                echo "/$filename" >> "$GITIGNORE_FILE"
                ((DIR_COUNT=DIR_COUNT+1))

            else
                echo "$filename" >> "$GITIGNORE_FILE"
                ((FILE_COUNT=FILE_COUNT+1))
            fi
            echo "${TICK} $filename"
            break
        done
    done    

    echo -e "\n${TICK}${NC}${RED}.gitignore${NC} Initialized in this directory."
    echo -e "${INFO}${NC}:"
    echo -e "\t${BOLD}Files Added${NC}: $FILE_COUNT \n\t${BOLD}Directories Added${NC}: $DIR_COUNT"
}

#!/bin/bash

# IMPORTS:
# Terminal Beauty
source "$RCOMMANDS_ROOT/configs/color.sh"
source "$RCOMMANDS_ROOT/configs/msg.sh"

help(){
    echo -e "\n ${BOLD} rgit || Ettisaf Rup ${NC} \n"
    echo "${TICK} ${RED}${BOLD}rgit${NC} — that's nothing but git!"
    echo -e "${INFO} here's help list for rgit \n"

    printf "%s\n" "    ${GREEN}${BOLD} > ${NC}rgit push
    ${SILVER} ${YELLOW}—${NC}  Automate initialization to pushing to the repository ${NC}
        ${RED}${BOLD} # ${NC}${SILVER}rgit push <commit-message>${NC}

    ${GREEN}${BOLD} > ${NC}rgit pull
    ${SILVER} ${YELLOW}—${NC}  Automate pulling before working contribution or set local branches ${NC}
         ${RED}${BOLD} # ${NC}${SILVER}rgit pull${NC}
                   
    ${GREEN}${BOLD} > ${NC}rgit rename
    ${SILVER} ${YELLOW}—${NC}  Simple command on renaming branch names ${NC}
         ${RED}${BOLD} # ${NC}${SILVER}rgit rename current_branch_name:new_branch_name${NC}
                   
    ${GREEN}${BOLD} > ${NC}rgit ignore
    ${SILVER} ${YELLOW}—${NC}  Creates an .gitignore file (way more simple) names${NC}
         ${RED}${BOLD} # ${NC}${SILVER}rgit ignore ${NC}"
}
#!/bin/bash

# Load colors and messages
source "$RCOMMANDS_ROOT/configs/color.sh"
source "$RCOMMANDS_ROOT/configs/msg.sh"

help() {
    echo -e "\n ${BOLD}${GREEN}rGit — That's nothing but git!${NC}"
    echo -e " ${TICK} ${YELLOW}${BOLD}Author:${NC} Ettisaf Rup"
    echo -e " ${INFO} Here's a list of available rgit commands:\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit push"
    echo -e "   ${SILVER}${YELLOW}—${NC} Automates init → commit → push"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit push \"your commit message\"${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit pull"
    echo -e "   ${SILVER}${YELLOW}—${NC} Pulls latest remote changes and updates local branch"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit pull${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit rename"
    echo -e "   ${SILVER}${YELLOW}—${NC} Renames current branch to a new one"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit rename current:new${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit cleanup"
    echo -e "   ${SILVER}${YELLOW}—${NC} Renames current branch to a new one"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit rename current:new${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit ignore"
    echo -e "   ${SILVER}${YELLOW}—${NC} Creates a smart .gitignore file"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit ignore${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit readme"
    echo -e "   ${SILVER}${YELLOW}—${NC} Generates a basic README.md file"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit readme${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit license"
    echo -e "   ${SILVER}${YELLOW}—${NC} Adds a LICENSE file (MIT by default)"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit license${NC}\n"

    echo -e " ${GREEN}${BOLD}>${NC} rgit archi"
    echo -e "   ${SILVER}${YELLOW}—${NC} Generates a project structure in a file named ARCHITECTURE"
    echo -e "   ${RED}${BOLD}#${NC} ${SILVER}rgit archi${NC}\n"

    echo -e " ${INFO} Tip: Run ${BOLD}rgit <command> --help${NC} for individual usage.\n"
}

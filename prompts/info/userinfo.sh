#!/bin/bash

 userinfo(){
        name=$(git config user.name)
        email=$(git config user.email)
        echo -e "${INFO} \n${BOLD}${YELLOW} User Git Identity: ${NC}"
        echo -e "\t${RED}# ${NC}${BOLD}User${NC}\t: ${name:-Not set}"
        echo -e "\t${RED}# ${NC}${BOLD}Email${NC}\t: ${email:-Not set}"
 }
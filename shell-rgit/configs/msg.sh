#!/bin/bash
# for Terminal Beauty (Colors, Bolding etc.)
source "$RCOMMANDS_ROOT/configs/color.sh"


# DATE GENERATOR
TODAY=$(date +"%d %B %Y")

# TICKS FOR TERMINAL MSGs
INFO="${BLUE}${BOLD}[INFO]${NC}"
TICK="${GREEN}${BOLD}✓${NC}"
CROSS="${RED}${BOLD}✗${NC}"

# REMOTE DEFINING MSGs
PROMPT_REMOTE_NAME="${BLUE}${BOLD}*${NC} Enter Remote name (e.g., origin):"
PROMPT_REPO_URL="${BLUE}${BOLD}*${NC} Enter Repository URL:"
# GIT Messages

NO_GIT_REMOTE_FOUND="${BOLD}No git remotes found. Kindly add one__ ${NC}"

# GIT Commit
COMMIT_EMPTY="${YELLOW}[rgit]: Commit message is empty! Add clean commit message:${NC}"

# SUCCESS/ERROR MESSAGES
SUCCESS="${GREEN}${BOLD}[rgit]: rgit executed successfully ${NC}"
FAIL="${GREEN}${BOLD}[rgit]: rgit executed with error, please check the faults. ${NC}"

# NO COMMAND FOUND
NO_COMMAND_FOUND="${CROSS} ${RED}${BOLD} No Command Found! ${NC}"
HELP_COMMAND="${INFO} ${BOLD} > Use 'rgit help' for details. ${NC}"
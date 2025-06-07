#!/bin/bash

license(){
LICENSE_DIR="$RCOMMANDS_ROOT/prompts/create/LICENSE"

echo -e "Choose ${BOLD}${RED}LICENSE${NC} for your repository:"
echo -e "${GREEN}${BOLD}#${NC} Available ${BOLD}${RED}LICENSES${NC}:"
echo "${BOLD}${YELLOW}1)${NC} MIT"
echo "${BOLD}${YELLOW}2)${NC} Apache-2.0"
echo "${BOLD}${YELLOW}3)${NC} GPL-3.0"
echo "${BOLD}${YELLOW}4)${NC} BSD-3-Clause"
echo "${BOLD}${YELLOW}5)${NC} Unlicense"
echo -n "${GREEN}${BOLD}>${NC} "
read choice

case $choice in
  1) license_file="$LICENSE_DIR/mit.txt" ;;
  2) license_file="$LICENSE_DIR/apache-2.0.txt" ;;
  3) license_file="$LICENSE_DIR/gpl-3.0.txt" ;;
  4) license_file="$LICENSE_DIR/bsd-3-clause.txt" ;;
  5) license_file="$LICENSE_DIR/unlicense.txt" ;;
  *) echo -e "${CROSS} Invalid choice" && exit 1 ;;
esac

if [ ! -f "$license_file" ]; then
  echo -e "${CROSS} ${RED}${BOLD}LICENSE${NC} file not found: $license_file"
  exit 1
fi

echo -en "\n📆 Enter year (default: $(date +%Y)): "
read year
year=${year:-$(date +%Y)}

echo -n "👤 Enter your name or organization: "
read fullname

# Copy and replace placeholders
cat "$license_file" | \
  sed "s/\[year\]/$year/g" | \
  sed "s/\[fullname\]/$fullname/g" > LICENSE

echo -e "${BOLD}\n\t${TICK} LICENSE file created successfully!"
}
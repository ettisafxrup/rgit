#!/bin/bash

archi() {

  OUTPUT_FILE="ARCHITECTURE"

  echo "📁 Application Structure:" > "$OUTPUT_FILE"
  echo >> "$OUTPUT_FILE"

  # Exclude hidden files, ARCHITECTURE file itself, and unnecessary folders
  find . \
    -not -path './.*' \
    -not -path './node_modules*' \
    -not -path './dist*' \
    -not -path './build*' \
    -not -path "./$OUTPUT_FILE" \
    -type f -o -type d | sort | while read path; do

    indent=$(echo "$path" | sed -e 's;[^/]*/;│   ;g' -e 's;│   \([^│]\);├── \1;')

    if [ -d "$path" ]; then
      echo "${indent/├── /├── 📂 }" >> "$OUTPUT_FILE"
    else
      echo "${indent/├── /├── 📑 }" >> "$OUTPUT_FILE"
    fi
  done

  echo -e "${TICK}${BOLD}${YELLOW} Application Structure ${NC}saved!\n\t ${GREEN}>${NC} ${SILVER}$OUTPUT_FILE${NC}"
}

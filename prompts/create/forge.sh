#!/bin/bash

forge(){

  # Output file
  OUTPUT_FILE="ARCHITECTURE"

  # Header
  echo "📁 Application Structure:" > "$OUTPUT_FILE"
  echo >> "$OUTPUT_FILE"

  # Tree generation with 📂 and 📑
  find . -not -name "$OUTPUT_FILE" | while read path; do
    indent=$(echo "$path" | sed -e 's;[^/]*/;│   ;g' -e 's;│   \([^│]\);├── \1;')
    if [ -d "$path" ]; then
      echo "${indent/├── /├── 📂 }" >> "$OUTPUT_FILE"
    else
      echo "${indent/├── /├── }" >> "$OUTPUT_FILE"
    fi
  done

  echo -e "${TICK}${BOLD}${YELLOW} Application Structure ${NC}saved!\n\t ${GREEN}>${NC} ${SILVER}$OUTPUT_FILE${NC}"
}

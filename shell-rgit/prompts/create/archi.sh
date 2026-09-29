archi() {
  OUTPUT_FILE="ARCHITECTURE"

  echo "📁 Application Structure:" > "$OUTPUT_FILE"
  echo >> "$OUTPUT_FILE"

  # Ignore patterns: hidden files/folders, build outputs, dependencies, IDEs, OS files, logs, etc.
  find . \
    \( \
      -path './.*' -a ! -path '.' -o \
      -name 'node_modules' -o \
      -name 'dist' -o \
      -name 'build' -o \
      -name 'bin' -o \
      -name 'obj' -o \
      -name '.vscode' -o \
      -name '.idea' -o \
      -name '.DS_Store' -o \
      -name 'Thumbs.db' -o \
      -name '.git' -o \
      -name '.svn' -o \
      -name '.hg' -o \
      -name '.env' -o \
      -name 'env' -o \
      -name 'venv' -o \
      -name '.venv' -o \
      -name '.mypy_cache' -o \
      -name '.pytest_cache' -o \
      -name '.coverage' -o \
      -name '.tox' -o \
      -name '.eggs' -o \
      -name 'bower_components' -o \
      -name 'coverage' -o \
      -name '.sass-cache' -o \
      -name '.cache' -o \
      -name 'outputs' -o \
      -name 'results' -o \
      -name 'logs' -o \
      -name '*.log' -o \
      -name '*.tmp' -o \
      -name '*.swp' -o \
      -name '*.swo' -o \
      -name '*.bak' -o \
      -name '*.apk' -o \
      -name '*.ipa' -o \
      -name '*.h5' -o \
      -name '*.ckpt' -o \
      -name '*.pth' -o \
      -name "$OUTPUT_FILE" \
    \) -prune -o -print | sort | while read path; do

    indent=$(echo "$path" | sed -e 's;[^/]*/;│   ;g' -e 's;│   \([^│]\);├── \1;')

    if [ -d "$path" ]; then
      echo "${indent/├── /├── 📂 }" >> "$OUTPUT_FILE"
    else
      echo "${indent/├── /├── }" >> "$OUTPUT_FILE"
    fi
  done

  echo -e "${TICK}${BOLD}${YELLOW} Application Structure ${NC}saved!\n\t ${GREEN}>${NC} ${SILVER}$OUTPUT_FILE${NC}"
}
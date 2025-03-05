#!/usr/bin/env bash

# Check if we're in a git repository
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "Error: Not in a git repository" >&2
    exit 1
fi


echo "File tree:"
# Use tree command with gitignore support if available (tree 2.0.0+)
tree --gitignore -a -I ".git|vendor|web|build"

echo -e "\nFile contents:"
echo "=============="

# Get all files tracked by git (respecting .gitignore)
git ls-files | grep -v "^vendor/" | grep -v "^web/" | grep -v "^internal/handlers/dist/" | while read -r file; do
    # Skip binary files
    if file "$file" | grep -q "binary"; then
        echo -e "\n$file (binary file, contents not shown)"
        echo "----------------------------------------"
    else
        echo -e "\n$file"
        echo "----------------------------------------"
        echo '```'
        cat "$file" 2>/dev/null || echo "Error: Could not read file"
        echo '```'
    fi
done

#!/bin/bash

login(){
        echo "🔐 Git Identity Setup"
        read -p "Git username: " name
        read -p "Git email: " email

        git config --global user.name "$name"
        git config --global user.email "$email"

        echo "✅ Identity configured."

        echo "Testing remote authentication..."
        git ls-remote &> /dev/null
        if [[ $? -eq 0 ]]; then
            echo "✅ Authenticated with remote."
        else
            echo "❌ Authentication failed. Try: git fetch or set up token."
        fi
        exit 0

}
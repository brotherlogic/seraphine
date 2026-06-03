#!/bin/bash

# Ensure the 'seraphine' session exists
if ! tmux has-session -t seraphine 2>/dev/null; then
  # Create a new session named 'seraphine', detached
  cd /workspaces/seraphine
  tmux new-session -d -s seraphine
 
  # Split the window horizontally (-h)
  # The left pane will remain a terminal
  # The right pane will run 'gh dash'
  tmux split-window -h -t seraphine "gh dash"
fi

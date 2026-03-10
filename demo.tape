# VHS tape file for generating the demo GIF
# Install VHS: https://github.com/charmbracelet/vhs
# Run: vhs demo.tape

Output demo.gif

Set FontSize 15
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms

Hide
Type "go run ."
Enter
Sleep 3s
Show

# Welcome screen — admire it
Sleep 2.5s

# Dismiss welcome → lands on lesson 1
Enter
Sleep 1.5s

# Go to lesson 2 (Context Window) — the star of the show
Type "n"
Sleep 1s

# Advance to visualization (bucket)
Right
Sleep 1.5s

# Add tools to fill the bucket
# Select "Read small file"
Enter
Sleep 0.6s

# Move to "Read large file" and add
Down
Enter
Sleep 0.6s

# Move to "Grep codebase" and add
Down
Down
Enter
Sleep 0.6s

# Move to "Bash: run tests" (big one) and add
Down
Down
Enter
Sleep 0.6s

# Add "Bash: build"
Down
Enter
Sleep 0.6s

# Add conversation messages
Down
Enter
Sleep 0.5s
Down
Enter
Sleep 0.5s

# Pause to show the filled bucket
Sleep 1.5s

# Switch to MCP servers tab
Type "2"
Sleep 1s

# Enable GitHub (34 tools!)
Enter
Sleep 0.8s

# Enable Slack
Down
Enter
Sleep 0.8s

# Enable Jira
Down
Enter
Sleep 0.8s

# Pause — show the MCP token cost
Sleep 2s

# Switch back to tools
Type "1"
Sleep 1s

# Add the huge file to trigger overflow + compression
Down
Down
Down
Down
Down
Down
Down
Down
Down
Enter
Sleep 2s

# Step through compression
Enter
Sleep 1.5s
Enter
Sleep 1.5s
Enter
Sleep 1.5s
Enter
Sleep 2s

# Return from compression
Enter
Sleep 1.5s

# Show the sidebar
Tab
Sleep 2s
Tab
Sleep 1s

# Jump to agentic loop lesson (lesson 15, intermediate)
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 0.5s
Type "n"
Sleep 1s

# Now on lesson 15 (Agentic Loop) — advance to viz
Right
Sleep 1.5s

# Step through the agentic loop iterations
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 0.8s
Enter
Sleep 2s

# Quit
Type "q"
Sleep 1s

# A seeded drill for the README: naming weekdays with the doomsday hint, then the breakdowns

Output assets/demo.gif

Set Shell bash
Set FontFamily "JetBrainsMono Nerd Font, JetBrains Mono, monospace"
Set FontSize 18
Set Width 900
Set Height 700
Set Padding 20
Set TypingSpeed 100ms
Set Theme "TokyoNight"

# Do setup invisibly
Hide
Type "./tuday --config Themes/tokyo-night.yaml --seed demo"
Enter
Sleep 500ms
Show

# Start capturing
Sleep 1s

# Turn on the doomsday hint
Type@300ms "jjjj"
Sleep 300ms
Enter
Sleep 800ms

# Start the drill
Type@200ms "kkkk"
Sleep 300ms
Enter
Sleep 2s

# Answer a few by number and by name, slipping up on the third
Type "1"
Sleep 2s
Type "5"
Sleep 2.5s
Type "4"
Sleep 2.5s
Type "wed"
Sleep 400ms
Enter
Sleep 2s
Type "4"
Sleep 2s
Type "fri"
Sleep 400ms
Enter
Sleep 2s

# Finish, and go through the breakdowns
Escape
Sleep 2.5s
Type "l"
Sleep 2s
Type "l"
Sleep 3s

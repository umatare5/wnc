# Configuration for https://github.com/charmbracelet/vhs/

Output wnc_show_overview_demo.gif

Set TypingSpeed 75ms
Set FontSize 14
Set Width 1400
Set Height 550

# Start of the tape

Type "ping -c 3 wnc1.example.internal"

Sleep 1s

Enter

Sleep 4s

Enter

Type "WNC_ACCESS_TOKEN=$(wnc generate token --username wnc --password password)"

Sleep 1s

Enter

Sleep 2s

Enter

Type "export WNC_CONTROLLERS=wnc1.example.internal:$WNC_ACCESS_TOKEN"

Sleep 1s

Enter

Sleep 2s

Enter

Type "wnc show overview --insecure"

Enter

Sleep 15s

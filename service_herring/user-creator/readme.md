# User Creator
Creates `n` amount of users with sudo powers, then delays for a random time and loops.
The users have a randomly selected name + a random number between 1-9999. 
## Parameters

`-v`: enable output

`n [num]`: Specify the number of users to create each loop

`--demo`: Display the generated usernames but don't create them

`-h` or `--help`: Display the help

Examples:
```
go run user-creator.go --demo -n 50

go build user-creator.go && ./user-creator -v -n 10
```
# burndown

Display the current state of a repository generated from different data on the project's github page.
The idea is to provide a tool that is simple to use to get a feel for if a team or project is
making productive changes.

##State

This project is currenlty in early stages, Some features that I plan to implement
- [ ] Make a 'guess' at how stale the project is
- [ ] Allow weights to be added to issues via labels
- [ ] Mark issues as ignored using labels
- [ ] Provide a display of issues fixed by recent commits
- [ ] Rank the most 'useful' committer (the one who fixes the most issues/the highest weighted issues)

Note: Nothing this project says means anything about your project, It is just using data to make a guess

## Installation

* Clone the repository into a valid go project path ie `git clone https://github.com/jlyon1/Burndown.git /home/jlyon1/burndown/src/burndown`

* Run `go get`

* Install Redis

* Launch redis with default configuration (For now)

* Launch with `go run main.go`

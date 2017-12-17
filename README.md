# Burndown [![Project State](https://burndown.io/badge/jlyon1/burndown)](https://burndown.io/jlyon1/burndown)

Display the current state of a repository generated from different data on the project's github page.
The idea is to provide a tool that is simple to use to get a feel for if a team or project is
making productive changes.

## About

How to use this project - The different charts represent different ways to look at your project in terms of issues, if you wish to ignore an issue add a 'issue/ignore' label, if you want to give it a different weight on the burn down chart (including zero and negative numbers) add a label burdown/{value}

I plan to add a more detailed README as time goes on, and implement smarter features for determining the state of the project.

The 'state' is currently determined by the time that the current open issues have been open, compared to the max time an issue has been open. If many issues have been open for a long time, the project is seen as more stale.

Badges are generated at /badge/{user}/{reponame}

## State of the project

This project is currenlty in early stages, Some features that I plan to implement
- [X] Make a 'guess' at how stale the project is
- [X] Allow weights to be added to issues via labels
- [X] Mark issues as ignored using labels
- [ ] Provide a display of issues fixed by recent commits
- [ ] Rank the most 'useful' committer (the one who fixes the most issues/the highest weighted issues)
- [ ] Add a config file
- [X] Use github api keys so more than 60 requests per hour can be made

Note: Nothing this project says means anything about your project, It is just using data to make a guess

## Installation

* Clone the repository into a valid go project path ie `git clone https://github.com/jlyon1/Burndown.git /home/jlyon1/burndown/src/burndown`

* Run `go get`

* Install Redis

* Launch redis with default configuration (For now)

* Launch with `go run main.go`

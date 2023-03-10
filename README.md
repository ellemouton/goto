# goto

Quickly navigate to a GitHub repo main page or specific commit from the command
line. 

### Note

Currently, this only works on OSx where `open` is installed by default.

## Install

```
go install github.com/ellemouton/goto@latest
```

## Recommended alias:
```
alias g2="goto go"
```

## Usage
```
// Register an alias for a new repo.
$ goto register <org> <repo> (optional <alias>)
$ goto register lightninglabs lnd

// Open a commit of a stored alias (assumes the g2 alias is set)
$ g2 lnd 74b9c9ce9

// Open a coomit of an unstored repo
$ g2 lightninglabs lnd 74b9c9ce9

// Just go to the repo main page
$ g2 lnd
```

## TODO

 - [ ] Add a `GitHost` to the `DBContext`. Have its default be `GitHub` but
   allow users to change the default to `GitLab` or w/e. Should still then be 
   able to register repos of other hosts. 
 - [ ] Add a `BrowserOpener` to the `DBContext` and set it depending on the
   OS that the binary is being run on _and_ let it be settable so that a user
   can override the default. ie: on first run, check OS & set accordingly and
   from then on, use whatever the DB has set. Add option for user to override. 
 - [ ] Add way to quickly navigate to the PR that added a commit. 
 - [ ] Cleanup code & add documentation
# goto

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

// Open a commit of a stored alias.
$ goto go lnd 74b9c9ce9

// Open a coomit of an unstored repo
$ goto go lightninglabs lnd 74b9c9ce9

// Just go to the repo main page
$goto go lnd
```
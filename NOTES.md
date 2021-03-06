# Notes

Potential new command reference? Each entity type has relevant commands, and then there
are shorter named commands for easy access to certain things.

What would the entry list command do? Could it be an easier way to provide entries for
completions without going through timesheet? It could use Bolt's looping over keys, via
some kind of abstraction.

```
$ tid entry [list] [--format=FORMAT]
$ tid entry create <DURATION> <NOTE> --timesheet=<DATE>
$ tid entry update <HASH> --offset=<OFFSET>
$ tid entry delete <HASH>

$ tid timesheet [list] [--format=FORMAT]
$ tid timesheet delete <DATE>

$ tid workspace [list] [--format=FORMAT]
$ tid workspace create <NAME>
$ tid workspace switch <NAME>
$ tid workspace update <NAME>
$ tid workspace delete <NAME>

$ tid start <NOTE>
$ tid resume [<HASH>]
$ tid status [<HASH>]
$ tid stop
$ tid report [--format=FORMAT] [--date=DATE] [--start=DATE] [--end=DATE]
```

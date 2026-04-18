# watchdiff

`watchdiff` is a tool, similar to `waatch` but instead of highligting changes, it prints a diff when ever there is a change in the output of the command being watched.

With `watchdiff`, you get the full initial command output followed by the history of all the changes as a sequence of unified diffs.

## Usage

```console
$ watchdiff -h
Watch a command and print colorized diffs of the output

Usage:
  watchdiff [command] [flags]

Flags:
  -C, --color string        Print colorized diffs (valid values are: auto, always or never) (default "auto")
  -c, --context int         Number of context lines for diff (default 4)
  -x, --exec                Run the command directly, not through the shell
  -h, --help                help for watchdiff
  -e, --include-stderr      Include stderr in the diff comparison (default true)
  -n, --interval duration   Interval between updates (e.g. 2s, 500ms) (default 2s)
  -q, --quiet               Suppress heartbeat dots
  -s, --shell string        Specify the shell to use (default "sh")
```


### Example output

```console
$ watchdiff -- ls -l
Monitoring: ls -l
Interval: 2s | Context: 4 | Stderr: true
---
INITIAL OUTPUT (21:34:10):
total 35
-rw-------. 1 u0_a329 u0_a329 8374 Apr 18 22:11 Makefile
-rw-------. 1 u0_a329 u0_a329 1047 Apr 18 22:33 README.md
drwx------. 2 u0_a329 u0_a329 3452 Apr 18 17:52 cmd
-rw-------. 1 u0_a329 u0_a329  248 Apr 18 17:01 go.mod
-rw-------. 1 u0_a329 u0_a329 1083 Apr 18 17:01 go.sum
-rw-------. 1 u0_a329 u0_a329  169 Apr 18 17:56 main.go
drwx------. 3 u0_a329 u0_a329 3452 Apr 18 17:49 pkg

---
.............
CHANGE DETECTED @ 21:34:38
@@ -5,4 +5,5 @@
 -rw-------. 1 u0_a329 u0_a329  248 Apr 18 17:01 go.mod
 -rw-------. 1 u0_a329 u0_a329 1083 Apr 18 17:01 go.sum
 -rw-------. 1 u0_a329 u0_a329  169 Apr 18 17:56 main.go
 drwx------. 3 u0_a329 u0_a329 3452 Apr 18 17:49 pkg
+-rw-------. 1 u0_a329 u0_a329    0 Apr 18 22:34 x.txt
........
CHANGE DETECTED @ 21:34:56
@@ -5,5 +5,4 @@
 -rw-------. 1 u0_a329 u0_a329  248 Apr 18 17:01 go.mod
 -rw-------. 1 u0_a329 u0_a329 1083 Apr 18 17:01 go.sum
 -rw-------. 1 u0_a329 u0_a329  169 Apr 18 17:56 main.go
 drwx------. 3 u0_a329 u0_a329 3452 Apr 18 17:49 pkg
--rw-------. 1 u0_a329 u0_a329    0 Apr 18 22:34 x.txt
....^C
signal: interrupt
```

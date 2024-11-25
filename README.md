Implementation based on video: https://www.youtube.com/watch?v=iSNsgj1OCLA

## Execution guide
Can be ran from the command line by:

```
go build
./prisoner-problem -n 100000 -p 100 -b 50
```

## Optional parameters
```
-n : Amount of executions to test (defaults to 1,000,000)
-p : Amount of prisoners to test against
-b : Amount of boxes a prisoner can check before exiting the room (must be <= prisoner count)
```

# god
![demo](./demo.gif)

## About
god is simple docker command history save tool.

## Require
- Docker

## Installtion
If you already installed golang, just go get.
```
go get -u github.com/skanehira/god
```

If you want download binary, please download from realese.

## Usage
```
# Select and run docker command from history
$ god

# Run docker command
# All arguments can be used as arguments are internally passed to the docker command
$ god ps -a
```

## Recommended settings
```
alias docker='/path/to/god'
```

## Storage location of history
```
# history storage is just file
$HOME/.docker_cmd_history
```

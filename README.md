# Amino Micro

A small microservice for `amino` encoding input data:

```bash
curl -XPOST localhost:3000/tx/encode --data-binary @tx.json
```

### Build and Run

To build and run `amino-micro` run:

```bash
$ make install
$ amino-micro


Usage:
  amino-micro [command]

Available Commands:
  help        Help about any command
  serve       Runs the server
  version     Prints version information

Flags:
      --config string   config file (default is $HOME/.amino-micro.yaml)
  -h, --help            help for amino-micro
  -t, --toggle          Help message for toggle

Use "amino-micro [command] --help" for more information about a command.
```

### Docker

To build the docker image with appropriate tags:

```bash
$ make docker
```

To push the docker image to the configured repo:

```bash
$ make docker-push
```

To run the just built docker image with the local config loaded and the proper port exposed:

```bash
$ make docker-run
```

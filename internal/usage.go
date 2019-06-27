package internal

const Usage = `cli, plugin-driven load testing framework.

Usage:

  cli [commands|flags]

The commands & flags are:

  version             print the version to stdout

  --config <file>                configuration file to load
  --debug                        turn on debug logging
  --pidfile <file>               file to write our pid to
  --quiet                        run in quiet mode
  --usage <plugin>               print usage for a plugin, ie, 'telegraf --usage mysql'
  --version                      display the version and exit

Examples:

  # run cli with all plugins defined in config file
  cli --config load.yaml

`

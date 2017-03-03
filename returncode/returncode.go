package returncode

// SUCCESS is the exit code in case of the a successful server stop
const SUCCESS = 0

// CONFIGFAILED is the exit code to use
// when the configuration failed to load
const CONFIGFAILED = 100

// HELP is the exit code triggered
// when user specify the help (-h or --help) flag
const HELP = 101

// CONFIGGEN is the exit code triggered
// when user specify the config generation (--config-gen) flag
const CONFIGGEN = 102

// ROUTERFAILED is the exit code to use
// when the router failed to start, or stopped brutally
const ROUTERFAILED = 110

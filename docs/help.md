# Help

The help text of every command but `completion`, transcribed from the binary.

## wnc

```plaintext
NAME:
   wnc - Operate Cisco Catalyst 9800 Wireless Network Controllers

USAGE:
   wnc [global options] [command [command options]]

VERSION:
   dev

COMMANDS:
   deauth             Deauthenticate a client on a controller
   delete             Delete a tag from a controller
   disable            Disable an access point or one of its radios
   enable             Enable an access point or one of its radios
   generate-token, g  Print the Basic auth token for a controller account
   reset              Restart an access point or its controller session
   save-config        Save the running configuration to the startup configuration
   set                Create or update a tag on a controller
   show, s            Display controller state

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
   --dry-run           report what would happen and change nothing
   --help, -h          show help
   --version, -v       print the version
```

## wnc deauth

```plaintext
NAME:
   wnc deauth - Deauthenticate a client on a controller

USAGE:
   wnc deauth (--mac <mac> | --username <username>) [options]

DESCRIPTION:
   --mac and --username are the mac and username columns of wnc show client,
   and one invocation gives one of them. The controller resolves it first, so a
   value it holds no client at is refused before the RPC, which answers the same
   whether or not a client was there. A username may hold more than one session,
   and the prompt says how many. The client is dropped and reconnects on its own
   within about four minutes. The operation is absent before 17.15. Pass --dry-run
   to name the target and change nothing.

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --mac string                                                       client MAC address, as shown in the mac column of wnc show client
   --username string                                                  client username, as shown in the username column of wnc show client
   --yes                                                              act without the confirmation prompt
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc delete

```plaintext
NAME:
   wnc delete - Delete a tag from a controller

USAGE:
   wnc delete [command [command options]]

COMMANDS:
   policy-tag  Delete one policy tag
   site-tag    Delete one site tag
   rf-tag      Delete one RF tag

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc delete policy-tag

```plaintext
NAME:
   wnc delete policy-tag - Delete one policy tag

USAGE:
   wnc delete policy-tag --name <name> [options]

DESCRIPTION:
   The name --name gives is read on the controller first, so a name it does not
   hold is a failure rather than a silent success. Pass --dry-run to report and
   change nothing.

OPTIONS:
   --name string  policy tag name, at most 32 characters
   --yes          act without the confirmation prompt
   --help, -h     show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc delete site-tag

```plaintext
NAME:
   wnc delete site-tag - Delete one site tag

USAGE:
   wnc delete site-tag --name <name> [options]

DESCRIPTION:
   The name --name gives is read on the controller first, so a name it does not
   hold is a failure rather than a silent success. Pass --dry-run to report and
   change nothing.

OPTIONS:
   --name string  site tag name, at most 32 characters
   --yes          act without the confirmation prompt
   --help, -h     show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc delete rf-tag

```plaintext
NAME:
   wnc delete rf-tag - Delete one RF tag

USAGE:
   wnc delete rf-tag --name <name> [options]

DESCRIPTION:
   The name --name gives is read on the controller first, so a name it does not
   hold is a failure rather than a silent success. Pass --dry-run to report and
   change nothing.

OPTIONS:
   --name string  RF tag name, at most 32 characters
   --yes          act without the confirmation prompt
   --help, -h     show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc disable

```plaintext
NAME:
   wnc disable - Disable an access point or one of its radios

USAGE:
   wnc disable [command [command options]]

COMMANDS:
   ap     Disable one access point
   radio  Disable one radio of one access point

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc disable ap

```plaintext
NAME:
   wnc disable ap - Disable one access point

USAGE:
   wnc disable ap --ap-name <ap-name> [options]

DESCRIPTION:
   --ap-name is the ap_name column of wnc show ap. The controller resolves it
   first, so a name it holds no access point under is refused before the RPC.
   This sets the access point's admin state, not one radio's.
   Pass --dry-run to name the target and change nothing.

OPTIONS:
   --ap-name string  access point name, as shown in the ap_name column of wnc show ap
   --yes             act without the confirmation prompt
   --help, -h        show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc disable radio

```plaintext
NAME:
   wnc disable radio - Disable one radio of one access point

USAGE:
   wnc disable radio --ap-name <ap-name> --slot <n> [options]

DESCRIPTION:
   --ap-name is the ap_name column of wnc show ap, and --slot is the Slot column
   of wnc show overview. The band the RPC needs is read from the controller and
   follows the radio type, so a dual-band radio takes one number either way.
   Pass --dry-run to name the target and change nothing.

OPTIONS:
   --ap-name string  access point name, as shown in the ap_name column of wnc show ap
   --yes             act without the confirmation prompt
   --slot int        radio slot, as shown in the Slot column of wnc show overview
   --help, -h        show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc enable

```plaintext
NAME:
   wnc enable - Enable an access point or one of its radios

USAGE:
   wnc enable [command [command options]]

COMMANDS:
   ap     Enable one access point
   radio  Enable one radio of one access point

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc enable ap

```plaintext
NAME:
   wnc enable ap - Enable one access point

USAGE:
   wnc enable ap --ap-name <ap-name> [options]

DESCRIPTION:
   --ap-name is the ap_name column of wnc show ap. The controller resolves it
   first, so a name it holds no access point under is refused before the RPC.
   This sets the access point's admin state, not one radio's.
   Pass --dry-run to name the target and change nothing.

OPTIONS:
   --ap-name string  access point name, as shown in the ap_name column of wnc show ap
   --yes             act without the confirmation prompt
   --help, -h        show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc enable radio

```plaintext
NAME:
   wnc enable radio - Enable one radio of one access point

USAGE:
   wnc enable radio --ap-name <ap-name> --slot <n> [options]

DESCRIPTION:
   --ap-name is the ap_name column of wnc show ap, and --slot is the Slot column
   of wnc show overview. The band the RPC needs is read from the controller and
   follows the radio type, so a dual-band radio takes one number either way.
   Pass --dry-run to name the target and change nothing.

OPTIONS:
   --ap-name string  access point name, as shown in the ap_name column of wnc show ap
   --yes             act without the confirmation prompt
   --slot int        radio slot, as shown in the Slot column of wnc show overview
   --help, -h        show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc generate-token

```plaintext
NAME:
   wnc generate-token - Print the Basic auth token for a controller account

USAGE:
   wnc generate-token [options]

OPTIONS:
   --username string, -u string  controller username [$WNC_USERNAME]
   --password string, -p string  controller password; prefer WNC_PASSWORD or piped stdin [$WNC_PASSWORD]
   --help, -h                    show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc reset

```plaintext
NAME:
   wnc reset - Restart an access point or its controller session

USAGE:
   wnc reset [command [command options]]

COMMANDS:
   ap      Restart one access point
   capwap  Reset one access point's controller session

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc reset ap

```plaintext
NAME:
   wnc reset ap - Restart one access point

USAGE:
   wnc reset ap --ap-name <ap-name> [options]

DESCRIPTION:
   --ap-name is the ap_name column of wnc show ap. The controller resolves it
   first, so a name it holds no access point under is refused before the RPC.
   Pass --dry-run to name the target and change nothing.

OPTIONS:
   --ap-name string  access point name, as shown in the ap_name column of wnc show ap
   --yes             act without the confirmation prompt
   --help, -h        show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc reset capwap

```plaintext
NAME:
   wnc reset capwap - Reset one access point's controller session

USAGE:
   wnc reset capwap --ap-name <ap-name> [options]

DESCRIPTION:
   --ap-name is the ap_name column of wnc show ap. The controller resolves it
   first, so a name it holds no access point under is refused before the RPC.
   The access point does not reboot: only its CAPWAP session is re-established.
   Pass --dry-run to name the target and change nothing.

OPTIONS:
   --ap-name string  access point name, as shown in the ap_name column of wnc show ap
   --yes             act without the confirmation prompt
   --help, -h        show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc save-config

```plaintext
NAME:
   wnc save-config - Save the running configuration to the startup configuration

USAGE:
   wnc save-config [options]

DESCRIPTION:
   The startup configuration is the only destination: no file may be named, and
   every change on the controller is persisted rather than only what this CLI
   wrote. An access point's admin state is unaffected, being no part of the
   configuration. The save took two to four seconds on every release measured, so
   a --timeout a read survives can still refuse it. Pass --dry-run to name the
   controller and change nothing.

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --yes                                                              act without the confirmation prompt
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc set

```plaintext
NAME:
   wnc set - Create or update a tag on a controller

USAGE:
   wnc set [command [command options]]

COMMANDS:
   policy-tag  Create or update one policy tag
   site-tag    Create or update one site tag
   rf-tag      Create or update one RF tag

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc set policy-tag

```plaintext
NAME:
   wnc set policy-tag - Create or update one policy tag

USAGE:
   wnc set policy-tag --name <name> [options]

DESCRIPTION:
   A name --name gives that the controller does not hold is created and one it
   holds is updated, so the same command may be repeated. A field no flag names
   is left as it is rather than cleared. Pass --dry-run to report and change
   nothing.

OPTIONS:
   --name string            policy tag name, at most 32 characters
   --yes                    act without the confirmation prompt
   --description string     description for the policy tag
   --wlan string            WLAN profile to bind, required with --policy-profile
   --policy-profile string  policy profile the WLAN is bound to, required with --wlan
   --help, -h               show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc set site-tag

```plaintext
NAME:
   wnc set site-tag - Create or update one site tag

USAGE:
   wnc set site-tag --name <name> [options]

DESCRIPTION:
   A name --name gives that the controller does not hold is created and one it
   holds is updated, so the same command may be repeated. A field no flag names
   is left as it is rather than cleared. Pass --dry-run to report and change
   nothing.

OPTIONS:
   --name string             site tag name, at most 32 characters
   --yes                     act without the confirmation prompt
   --description string      description for the site tag
   --ap-join-profile string  AP join profile to bind
   --flex-profile string     flex profile to bind, which clears --local-site
   --local-site              mark the site local
   --help, -h                show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc set rf-tag

```plaintext
NAME:
   wnc set rf-tag - Create or update one RF tag

USAGE:
   wnc set rf-tag --name <name> [options]

DESCRIPTION:
   A name --name gives that the controller does not hold is created and one it
   holds is updated, so the same command may be repeated. A field no flag names
   is left as it is rather than cleared. Pass --dry-run to report and change
   nothing.

OPTIONS:
   --name string           RF tag name, at most 32 characters
   --yes                   act without the confirmation prompt
   --description string    description for the RF tag
   --profile-24ghz string  2.4 GHz RF profile to bind
   --profile-5ghz string   5 GHz RF profile to bind
   --profile-6ghz string   6 GHz RF profile to bind
   --help, -h              show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port] [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for the controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
```

## wnc show

```plaintext
NAME:
   wnc show - Display controller state

USAGE:
   wnc show [command [command options]]

COMMANDS:
   overview, o    Per-radio RF summary across 2.4, 5 and 6 GHz
   ap, a          Associated access points
   ap-join, join  Join, discovery and DTLS outcome per access point, joined or not
   ap-tag, tag    Tag assignment and its resolved values, per access point
   client, c      Associated wireless clients
   wlan, w        Configured WLANs and their bound policy profiles
   policy-tag     Configured policy tags and the WLANs they bind
   site-tag       Configured site tags and the profiles they name
   rf-tag         Configured RF tags and their per-band RF profiles

OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
   --help, -h                                                         show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")
```

## wnc show overview

```plaintext
NAME:
   wnc show overview - Per-radio RF summary across 2.4, 5 and 6 GHz

USAGE:
   wnc show overview [options]

DESCRIPTION:
   One row per access point radio, sorted by ap_name.
   A cell reading "-" is a value the controller did not send.
   Admin is the radio's own state: an access-point-level disable leaves it
   Enabled with Oper reading Down, and wnc show ap is the authority instead.

OPTIONS:
   --sort-by key, -b key      sort key (see --sort-keys) (default: "ap_name")
   --sort-keys                print the keys --sort-by accepts, then exit
   --radio string, -r string  band filter (2.4|5|6)
   --help, -h                 show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show ap

```plaintext
NAME:
   wnc show ap - Associated access points

USAGE:
   wnc show ap [options]

DESCRIPTION:
   One row per access point in capwap-data, sorted by ap_name.
   A cell reading "-" is a value the controller did not send.
   Admin is the access point's own state, which an access-point-level disable
   changes and the Admin column of wnc show overview does not.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "ap_name")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show ap-join

```plaintext
NAME:
   wnc show ap-join - Join, discovery and DTLS outcome per access point, joined or not

USAGE:
   wnc show ap-join [options]

DESCRIPTION:
   One row per access point the controller remembers, joined or not, sorted
   by ap_name.
   A cell reading "-" is a value the controller did not send.
   It is the only view that reports an access point capwap-data has dropped.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "ap_name")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show ap-tag

```plaintext
NAME:
   wnc show ap-tag - Tag assignment and its resolved values, per access point

USAGE:
   wnc show ap-tag [options]

DESCRIPTION:
   One row per access point, sorted by ap_name.
   A cell reading "-" is a value the controller did not send.
   The tag columns are the resolved tags in force; the two profile columns
   come from the configured site tag and agree only while Tag Source is Static.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "ap_name")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show client

```plaintext
NAME:
   wnc show client - Associated wireless clients

USAGE:
   wnc show client [options]

DESCRIPTION:
   One row per associated client, sorted by mac.
   A cell reading "-" is a value the controller did not send.
   --radio, --ssid and --ap-name narrow the list. A client whose band the
   controller did not report is excluded by --radio, and the count is logged.

OPTIONS:
   --sort-by key, -b key      sort key (see --sort-keys) (default: "mac")
   --sort-keys                print the keys --sort-by accepts, then exit
   --radio string, -r string  band filter (2.4|5|6)
   --ssid string, -s string   keep only clients on this SSID
   --ap-name string           keep only clients on this AP
   --help, -h                 show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show wlan

```plaintext
NAME:
   wnc show wlan - Configured WLANs and their bound policy profiles

USAGE:
   wnc show wlan [options]

DESCRIPTION:
   One row per WLAN and each policy profile bound to it, sorted by wlan_id,
   so a WLAN bound under two tags appears twice and an unbound one appears
   once with its policy columns empty.
   A cell reading "-" is a value the controller did not send.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "wlan_id")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show policy-tag

```plaintext
NAME:
   wnc show policy-tag - Configured policy tags and the WLANs they bind

USAGE:
   wnc show policy-tag [options]

DESCRIPTION:
   One row per WLAN binding the tag carries, sorted by policy_tag, so a tag
   binding three WLANs appears three times and one binding none appears once.
   A cell reading "-" is a value the controller did not send.
   WLAN is the WLAN profile name the binding keys on, not always the SSID.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "policy_tag")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show site-tag

```plaintext
NAME:
   wnc show site-tag - Configured site tags and the profiles they name

USAGE:
   wnc show site-tag [options]

DESCRIPTION:
   One row per site tag, sorted by site_tag.
   A cell reading "-" is a value the controller did not send.
   The read asks for the values in force, so a leaf a tag left at its default
   is reported rather than arriving as an absence.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "site_tag")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

## wnc show rf-tag

```plaintext
NAME:
   wnc show rf-tag - Configured RF tags and their per-band RF profiles

USAGE:
   wnc show rf-tag [options]

DESCRIPTION:
   One row per RF tag, sorted by rf_tag.
   A cell reading "-" is a value the controller did not send.
   The read asks for the values in force, because a plain read omits the
   built-in tag's three per-band profile names.

OPTIONS:
   --sort-by key, -b key  sort key (see --sort-keys) (default: "rf_tag")
   --sort-keys            print the keys --sort-by accepts, then exit
   --help, -h             show help

GLOBAL OPTIONS:
   --config string     path to the JSON configuration file [$WNC_CONFIG]
   --log-level string  log verbosity (error|warning|debug) (default: "warning")

INHERITED OPTIONS:
   --controller string, -c string [ --controller string, -c string ]  controller host[:port], repeatable [$WNC_CONTROLLER]
   --access-token string                                              Basic auth token for every controller [$WNC_ACCESS_TOKEN]
   --insecure, -k                                                     skip TLS certificate verification
   --format string, -o string, -f string                              output format (table|json) (default: "table")
   --pretty                                                           draw the table with borders and status glyphs
   --timeout duration, -t duration                                    request timeout (default: 1m0s)
   --sort-order string                                                sort direction (asc|desc) (default: "asc")
```

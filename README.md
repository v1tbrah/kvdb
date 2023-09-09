# Key Value Database
### Restrictions
* Supports only strings as key, and only strings as value
### Start
* Run daemon
```shell
go run main.go
```
* Connect to daemon
```shell
telnet localhost 4321
```
* Set something
```SET foo bar```
* Get something
```GET foo```
* Delete something
```DELETE foo```
### Config
* See [there](config/config.go)
### Syntax
* SET operation:    SET KEY VALUE
* GET operation:    GET KEY
* DELETE operation: DELETE KEY

# Primus

## Usage

### Server

```
$ primus-server -c config.toml
```

```toml:config.toml
# Configuration

[server]
port = 14300

[log]
level = "debug"
# access_log = "access_log"
# error_log = "error_log"
```

### Client

```
$ primus-client -c config.toml
```

```toml:config.toml
[server]
host = "localhost"
port = 14300
ssl = false

[[Route]]
channel = "sample1"

host = "localhost"
port = 8080
path = "/sample/endpoint"
ssl  = false
```


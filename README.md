IX2005 Exporter
---------------

Export NEC IX2005 status.

This exporter scrape data from web ui so you can use this when you cannot use SNMP. If you can use SNMP, you should use [SNMP Exporter](https://github.com/prometheus/snmp_exporter) instead.

## usage

``` shell
$ make
$ ./ix2005-exporter [flags]
```

### flags

``` shell
$ ./ix2005-exporter --help
```

* `ix2005.uri`: URI of target IX2005.
* `web.listen-address`: Address to listen on for web interface and telemetry.
* `web.telemetry-path`: Path under which to expose metrics.
* `version`: Show application version.

### exported metrics

|          metric           |                  meaning                   | label  |
|:-------------------------:|:------------------------------------------:|:------:|
| ix2005_inside_temperature | The temperature of inside of target IX2005 | target |
| ix2005_memory             | Memory usage of target IX2005              | target |

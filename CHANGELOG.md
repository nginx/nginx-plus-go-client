# Changelog

Starting with version 0.8.0, an automatically generated list of changes can be found on the [GitHub Releases page](https://github.com/nginxinc/nginx-plus-go-client/releases).

## 0.7.0 (Jul 10, 2020)

FEATURES:

- [38](https://github.com/nginxinc/nginx-plus-go-client/pull/38): *Support for /slabs API endpoint*. The client now
  supports retrieving shared memory zone usage info.
- [41](https://github.com/nginxinc/nginx-plus-go-client/pull/41): *Support for /processes API endpoint*. The client now
  supports retrieving processes info.

CHANGES:

- The version of NGINX Plus for e2e testing was changed to R22.
- The version of Go was changed to 1.14

## 0.6.0 (Nov 8, 2019)

FEATURES:

- [34](https://github.com/nginxinc/nginx-plus-go-client/pull/34): *Support for updating upstream servers parameters*.
  The client now supports updating upstream parameters of servers that already exist in NGINX Plus.

CHANGES:

- Public methods `UpdateHTTPServers` and `UpdateStreamServers` now return a third slice that includes the updated
  servers -- i.e. the servers that were already present in NGINX Plus but were updated with different parameters.
- Client will assume port `80` in addresses of updated servers of `UpdateHTTPServers` and `UpdateStreamServers` if port
  is not explicitly set.
- The version of Go was changed to 1.13

## 0.5.0 (Sep 25, 2019)

FEATURES:

- [30](https://github.com/nginxinc/nginx-plus-go-client/pull/30): *Support additional upstream server parameters*. The
client now supports configuring `route`, `backup`, `down`, `drain`,  `weight` and `service` parameters for http
upstreams and  `backup`, `down`,  `weight` and  `service` parameters for stream upstreams.
- [31](https://github.com/nginxinc/nginx-plus-go-client/pull/31): *Support location zones and resolver metrics*.

FIXES:

- [29](https://github.com/nginxinc/nginx-plus-go-client/pull/29): *Fix max_fails parameter in upstream servers*.
  Previously, if the MaxFails field was not explicitly set, the client would incorrectly configure an upstream with the
  value `0` instead of the correct value `1`.

CHANGES:

- The version of NGINX Plus for e2e testing was changed to R19.
- The version of the API was changed to 5.

## 0.4.0 (July 17, 2019)

FEATURES:

- [24](https://github.com/nginxinc/nginx-plus-go-client/pull/24): *Support `MaxConns` in upstream servers*.

FIXES:

- [25](https://github.com/nginxinc/nginx-plus-go-client/pull/25): *Fix session metrics for stream server zones*. Session
  metrics with a status of `4xx` or `5xx` are now correctly reported. Previously they were always reported as `0`.

## 0.3.1 (June 10, 2019)

CHANGES:

- [22](https://github.com/nginxinc/nginx-plus-go-client/pull/22): *Change in stream zone sync metrics*. `StreamZoneSync`
  field of the `Stats` type is now a pointer. It will be nil if NGINX Plus doesn't report any zone sync stats.

## 0.3 (May 29, 2019)

FEATURES:

- [20](https://github.com/nginxinc/nginx-plus-go-client/pull/20): *Support for stream zone sync metrics*. The client
  `GetStats` method now additionally returns stream zone sync metrics.
- [13](https://github.com/nginxinc/nginx-plus-go-client/pull/13): *Support for key-value endpoints*. The client
  implements a set of methods to create/modify/delete key-val pairs for both http and stream contexts.
- [12](https://github.com/nginxinc/nginx-plus-go-client/pull/12) *Support for NGINX status info*. The client `GetStats`
  method now additionally returns NGINX status metrics. Thanks to [jthurman42](https://github.com/jthurman42).

CHANGES:

- The repository was renamed to `nginx-plus-go-client` instead of `nginx-plus-go-sdk`. If the client is used as a
  dependency, this name needs to be changed in the import section (`import
  "github.com/nginxinc/nginx-plus-go-client/client"`).
- The version of the API was changed to 4.
- The version of NGINX Plus for e2e testing was changed to R18.

## 0.2 (Sep 7, 2018)

FEATURES:

- [7](https://github.com/nginxinc/nginx-plus-go-sdk/pull/7): *Support for stream server zone and stream upstream
  metrics*. The client `GetStats` method now additionally returns stream server zone and stream upstream metrics.

CHANGES:

- The version of NGINX Plus for e2e testing was changed to R16.

## 0.1 (July 30, 2018)

Initial release

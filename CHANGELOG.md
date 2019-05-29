## 0.3 (May 29, 2019)
FEATURES:
* [20](https://github.com/nginxinc/nginx-plus-go-client/pull/20): *Support for stream zone sync metrics*. The client `GetStats` method now additionally returns stream zone sync metrics. 
* [13](https://github.com/nginxinc/nginx-plus-go-client/pull/13): *Support for key-value endpoints*. The client implements a set of methods to create/modify/delete key-val pairs for both http and stream contexts.
* [12](https://github.com/nginxinc/nginx-plus-go-client/pull/12) *Support for NGINX status info*. The client `GetStats` method now additionally returns NGINX status metrics. Thanks to [jthurman42](https://github.com/jthurman42).

CHANGES:
* The repository was renamed to `nginx-plus-go-client` instead of `nginx-plus-go-sdk`. If the client is used as a dependency, this name needs to be changed in the import section (`import "github.com/nginxinc/nginx-plus-go-client/client"`).
* The version of the API was changed to 4.
* The version of NGINX Plus for e2e testing was changed to R18.

## 0.2 (Sep 7, 2018)

FEATURES:
* [7](https://github.com/nginxinc/nginx-plus-go-sdk/pull/7): *Support for stream server zone and stream upstream metrics*. The client `GetStats` method now additionally returns stream server zone and stream upstream metrics.

CHANGES:
* The version of NGINX Plus for e2e testing was changed to R16.

## 0.1 (July 30, 2018)
Initial release

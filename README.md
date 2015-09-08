# simplehttps

[![Coverage Status](https://coveralls.io/repos/nfisher/simplehttps/badge.svg?branch=master&service=github)](https://coveralls.io/github/nfisher/simplehttps?branch=master)
[![Build Status](https://travis-ci.org/nfisher/simplehttps.svg?branch=master)](https://travis-ci.org/nfisher/simplehttps)


Simple HTTPS server was developed to simplify development for secure sites. It aims to provide a simple mechanism to proxy APIs and a static app under the same domain on a developers workstation. Thereby eliminating the need for CORS exceptions and encouraging the use of best practises like HTTPS only cookies.

## Dependencies

- self-signed certificate.
- static site.
- IP based applications (virtual-host mappings currently not supported).
- Golang >= 1.4.2.

## Getting Started

First up you'll need a JSON configuration file. You can modify the config.json that's found in the root of the respository a sample configuration is as follows;

```
{
  "apps": {
    "/cms": "http://127.0.0.1:8080",
    "//dev.local:8443/cms": "http://127.0.0.1:8081",
  }
}
```

With the above configuration;

- requests to https://dev.local:8443/app1/* will be routed to the application listening at 127.0.0.1:8080.
- requests to /cms/* with any host but "dev.local:8443" will be routed to the application listening at 127.0.0.1:8080.
- any requests that do not match the above two criteria will be served from the static folder (default _site).

The easiest way to route to the dev.local domain is to add a host entry. Alternatively you can set the host header if using curl or equivalent.

# Listen on 443

Port 443 is considered a privileged port and most if not all systems requires administrative privileges to be bound. However if you want clean https URLs such as https://dev.local/ without specifying a port you'll need to start the service under a user account with administrative privileges.

Given this app is developer focused I would suggest binding to the loop back interface as follows as an administrator;

```
simplehttps -listen="127.0.0.1:443"
```


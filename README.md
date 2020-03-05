# gnucash-graphsql

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/vinymeuh/gnucash-graphql.svg)](https://github.com/vinymeuh/gnucash-graphql/releases/latest)
[![Build Status](https://travis-ci.org/vinymeuh/gnucash-graphql.svg?branch=master)](https://travis-ci.org/vinymeuh/gnucash-graphql)
[![codecov](https://codecov.io/gh/vinymeuh/gnucash-graphql/branch/master/graph/badge.svg)](https://codecov.io/gh/vinymeuh/gnucash-graphql)
[![Go Report Card](https://goreportcard.com/badge/github.com/vinymeuh/gnucash-graphql)](https://goreportcard.com/report/github.com/vinymeuh/gnucash-graphql)

## How to start the server

By default configuration file ```config.yml``` is loaded from current directory but you can use another using environment variable  ```GNCGQL_CONF```. See ```config.yml.template``` for more details.

```
~> export GNCGQL_CONF=/path/to/my/config.yml
~> ./gnucash-graphql
```


When is started, we can get some informations about the server using command:

```
~> curl -s localhost:3000 | python -m json.tool
```

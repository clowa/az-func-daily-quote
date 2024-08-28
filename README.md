[![Tests](https://github.com/clowa/az-func-daily-quote/actions/workflows/run-tests.yaml/badge.svg)](https://github.com/clowa/az-func-daily-quote/actions/workflows/run-tests.yaml)

# Daily quote API

This repository contains the infrastructure and code of an Azure Function App that serves an API for quotes.

## Repository structure

- [`function`](./function/README.md): Contains the code of the Azure Function App.
- `terraform`: Contains the infrastructure definition of the Azure Function App.
- [`docs`](./docs/README.md): Contains general documentation about this project. Please see the `README.md` of each specific function for more detailed information.

## Docs

Documentation can be found in the [`docs`](./docs/README.md) folder.

## Facts

- Average cold start time of `quotes`: ~ 4396 ms
- Average hot execution time of `quotes`: ~ 291 ms
- [Hey](https://github.com/rakyll/hey) stats:

  ```bash
  $ hey -m GET https://quotes.clowa.dev/api/quote

  Summary:
  Total:        38.4224 secs
  Slowest:      18.2495 secs
  Fastest:      4.1377 secs
  Average:      7.7343 secs
  Requests/sec: 5.2053

  Total data:   45597 bytes
  Size/request: 229 bytes

    Response time histogram:
    4.138 [1] |
    5.549 [52] |■■■■■■■■■■■■■■■■■■■■■■
    6.960 [0] |
    8.371 [96] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
    9.782 [0] |
    11.194 [0] |
    12.605 [48] |■■■■■■■■■■■■■■■■■■■■
    14.016 [0] |
    15.427 [1] |
    16.838 [0] |
    18.249 [1] |

    Latency distribution:
    10% in 4.3543 secs
    25% in 4.4647 secs
    50% in 7.2960 secs
    75% in 11.7991 secs
    90% in 11.9979 secs
    95% in 12.0300 secs
    99% in 18.2495 secs

    Details (average, fastest, slowest):
    DNS+dialup: 0.0621 secs, 4.1377 secs, 18.2495 secs
    DNS-lookup: 0.0098 secs, 0.0000 secs, 0.0399 secs
    req write: 0.0000 secs, 0.0000 secs, 0.0001 secs
    resp wait: 7.6718 secs, 4.1376 secs, 18.0013 secs
    resp read: 0.0001 secs, 0.0000 secs, 0.0007 secs

    Status code distribution:
    [200] 199 responses
  ```

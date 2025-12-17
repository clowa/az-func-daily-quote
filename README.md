[![Tests](https://github.com/clowa/az-func-daily-quote/actions/workflows/run-tests.yaml/badge.svg)](https://github.com/clowa/az-func-daily-quote/actions/workflows/run-tests.yaml)

# Daily quote API

This repository contains the infrastructure and code of an Azure Function App that serves an API for quotes.

## Repository structure

- [`function`](./function/README.md): Contains the code of the Azure Function App.
- [`terraform`](./terraform/README.md): Contains the infrastructure definition of the Azure Function App.
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
    Total:        6.5917 secs
    Slowest:      5.9453 secs
    Fastest:      0.1018 secs
    Average:      0.6637 secs
    Requests/sec: 30.3410
    
    Total data:   41000 bytes
    Size/request: 205 bytes

  Response time histogram:
    0.102 [1]     |
    0.686 [147]   |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
    1.271 [11]    |■■■
    1.855 [24]    |■■■■■■■
    2.439 [1]     |
    3.024 [13]    |■■■■
    3.608 [2]     |■
    4.192 [0]     |
    4.777 [0]     |
    5.361 [0]     |
    5.945 [1]     |


  Latency distribution:
    10% in 0.1321 secs
    25% in 0.1594 secs
    50% in 0.2466 secs
    75% in 0.7190 secs
    90% in 1.8299 secs
    95% in 2.7432 secs
    99% in 3.2167 secs

  Details (average, fastest, slowest):
    DNS+dialup:   0.0595 secs, 0.1018 secs, 5.9453 secs
    DNS-lookup:   0.0211 secs, 0.0000 secs, 0.0851 secs
    req write:    0.0000 secs, 0.0000 secs, 0.0004 secs
    resp wait:    0.6040 secs, 0.1017 secs, 5.7147 secs
    resp read:    0.0001 secs, 0.0000 secs, 0.0032 secs

  Status code distribution:
    [200] 200 responses
  ```

<p align="center">
  <img src="docs/bbq.png" width="200" alt="bbq logo" />
</p>

<h1 align="center">bbq</h1>

<p align="center">
  <strong>a collection of algorithms and data structures in go</strong>
</p>

<div align="center">

<br>

[![go version][go_version_img]][go_dev_url]
[![version][repo_version_img]][repo_version_url]
[![license][repo_license_img]][repo_license_url]

<br>

</div>

---

`bbq` is a small collection of useful data structures and algorithms reused across projects

```sh
go get kamjapowered.com/bbq
```

the import path is `kamjapowered.com/bbq`, not the github url. the repository is
the source; the module path is the domain.

| package       | description                                                                                 |
| ------------- | ------------------------------------------------------------------------------------------- |
| `constraints` | generic numeric and ordering constraints used across the module                             |
| `control`     | reference signals and simple control-system helpers for time-based values                   |
| `csv`         | small csv helpers for reading and working with tabular data                                 |
| `curves`      | interpolation and easing curves for normalized progress values                              |
| `errors`      | typed errors with optional fields, wrapping, and terminal-aware formatting                  |
| `filter`      | online filters for smoothing, limiting, and estimating scalar values                        |
| `functional`  | generic helpers for mapping values and applying functional options                          |
| `geometry`    | generic 2d geometry primitives built around vectors, shapes, paths, and rays                |
| `hashmap`     | hash-based structures including dense maps, bidirectional maps, and typemaps                |
| `linked_list` | singly, doubly, and circular linked list implementations                                    |
| `math`        | generic numeric helpers for common arithmetic, interpolation, and statistics                |
| `queue`       | fifo, ring, bounded, double-ended, linked-list, and priority queue implementations          |
| `set`         | set-like structures including hash sets, bit sets, checklists, and generational populations |
| `sort`        | generic sorting algorithms for ordered values                                               |
| `stack`       | slice-backed, linked-list, and unique lifo stack implementations                            |
| `strings`     | string distance metrics and fuzzy matching helpers                                          |
| `system`      | small wrappers for process and operating-system interactions                                |
| `terminal`    | ansi escape constants and terminal capability helpers                                       |

## licence

the code is [mit](./LICENSE).

**the content is carved out.** the name, the logo and the wordmark are all
rights reserved. the mit grant covers the code in this repository and not the
marks it is published under.

[go_version_img]: https://img.shields.io/badge/go-1.27+-00add8?style=for-the-badge&logo=go
[go_dev_url]: https://go.dev/
[repo_license_img]: https://img.shields.io/github/license/kamjapowered/bbq?style=for-the-badge
[repo_license_url]: https://github.com/kamjapowered/bbq/blob/main/LICENSE
[repo_version_img]: https://img.shields.io/github/v/tag/kamjapowered/bbq?style=for-the-badge
[repo_version_url]: https://github.com/kamjapowered/bbq/tags

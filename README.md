<a name="readme-top"></a>

<div align="center">

# go-dsa  
### algorithms and data structures in go

a curated collection of practical algorithms and data structures i reuse across my go projects

<br>

[![go version][go_version_img]][go_dev_url]
[![version](https://img.shields.io/github/v/tag/vistormu/go-dsa?style=for-the-badge)](https://github.com/vistormu/go-dsa/tags)
[![License][repo_license_img]][repo_license_url]

<br>

</div>

> [!WARNING]  
> this is a personal project under active development  
> expect breaking changes, slow iteration and occasional bugs

---

## philosophy

this repository is **not** an academic catalogue.

it focuses on:
- practical data structures
- small, reusable algorithms
- minimal apis
- allocation-aware implementations
- utilities missing or awkward in the standard library

everything here exists because it has been needed in real projects  
(robotics, simulation, engines, tooling, servers)

---

## table of contents

- [buffer](#buffer)
- [constraints](#constraints)
- [control](#control)
- [csv](#csv)
- [errors](#errors)
- [filter](#filter)
- [geometry](#geometry)
- [hashmap](#hashmap)
- [linked_list](#linked_list)
- [math](#math)
- [queue](#queue)
- [set](#set)
- [sort](#sort)
- [stack](#stack)
- [strings](#strings)
- [system](#system)
- [terminal](#terminal)

---

## buffer

low level byte and ring buffers used for io, streaming and internal queues

---

## constraints

generic numeric and type constraints used across the repository

this package exists to avoid repeating constraint definitions and to keep generics readable

---

## control

control theory primitives and reference generators

includes:
- pid controller with derivative filtering and anti-windup
- first and second order systems (state space, discrete)
- reference generators: step, ramp, sine, square, triangular

designed for simulation, robotics and real-time systems

---

## csv

small helpers for csv parsing and generation

intentionally minimal and allocation aware

---

## errors

structured error values with:
- typed error kinds
- key-value metadata
- error wrapping
- pretty formatting via `fmt`

this package does **not** replace `error`  
it enhances error inspection and presentation

---

## filter

time series filters and signal conditioning blocks

includes:
- mean and median filters
- low pass filter
- rate limiter
- dead zone
- kalman filter (scalar)
- kalman filter with constant velocity model

designed for online, incremental use

---

## geometry

generic geometric primitives built around a single `Vector` type

includes:
- vector (with optional z)
- segment, line, ray
- rect, capsule, ellipse, polygon
- paths and arrows

no position ownership or transforms  
pure geometry only

---

## hashmap

hash based data structures built on top of go maps

includes:
- bidirectional hashmap
- typemap keyed by reflect.Type

focused on correctness and clarity over cleverness

---

## linked_list

linked list variants for educational and niche use cases

includes:
- singly linked list
- doubly linked list
- circular linked list

most code should prefer slices, but these are useful when needed

---

## math

generic numeric helpers missing from go’s standard library

includes:
- abs, sign, clamp, remap
- lerp and inverse lerp
- sum, mean, variance
- min/max over slices
- floating point comparisons

designed to work with generics

---

## queue

fifo data structures with different tradeoffs

includes:
- array backed queue
- ring buffer queue
- bounded queue
- linked list queue
- deque
- priority queue

all implementations avoid panics and expose simple apis

---

## set

set-like data structures

includes:
- hash set
- bit set
- checklist (completion tracking)

useful for ecs, scheduling and state tracking

---

## sort

sorting algorithms for learning and controlled use

includes:
- quicksort
- tests and benchmarks

for most code, prefer `slices.Sort`

---

## stack

lifo data structures

includes:
- array stack
- linked list stack
- unique stack (deduplicated, order preserving)

useful for parsers, undo stacks and graph traversal

---

## strings

string similarity and distance metrics

includes:
- levenshtein
- damerau-levenshtein
- jaro and jaro-winkler
- hamming
- trigram jaccard
- fuzzy search helpers

useful for search, ranking and matching

---

## system

small os and runtime helpers

includes:
- signal listener for graceful shutdown

keeps system concerns out of application logic

---

## terminal

ansi and terminal helpers

includes:
- colour codes
- window utilities

used by the errors package and cli tools

---

## contributing

this repository is primarily maintained for personal use  
however, issues and pull requests are welcome if they:

- keep apis minimal
- avoid unnecessary abstraction
- do not duplicate standard library features

no commitment is made to long term api stability

---

[go_version_img]: https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go
[go_dev_url]: https://go.dev/
[repo_license_img]: https://img.shields.io/github/license/vistormu/go-dsa?style=for-the-badge
[repo_license_url]: https://github.com/vistormu/go-dsa/blob/main/LICENSE

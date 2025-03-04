---
layout: post
title: ADVANCED
permalink: /docs/cli/advanced
redirect_from:
 - /cli/advanced.md/
 - /docs/cli/advanced.md/
---

# `ais advanced` commands

Commands for special use cases (e.g. scripting) and *advanced* usage scenarios, whereby a certain level of understanding of possible consequences is implied and required:

```console
$ ais advanced --help
NAME:
   ais advanced - special commands intended for development and advanced usage

USAGE:
   ais advanced command [command options] [arguments...]

COMMANDS:
   gen-shards        generate and write random TAR shards, e.g.:
                     - gen-shards 'ais://bucket1/shard-{001..999}.tar' - write 999 random shards (default sizes) to ais://bucket1
                     - gen-shards 'gs://bucket2/shard-{01..20..2}.tgz' - 10 random gzipped tarfiles to Cloud bucket
                     (notice quotation marks in both cases)
   resilver          resilver user data on a given target (or all targets in the cluster): fix data redundancy
                     with respect to bucket configuration, remove migrated objects and old/obsolete workfiles
   preload           preload object metadata into in-memory cache
   remove-from-smap  immediately remove node from cluster map (advanced usage - potential data loss!)
   random-node       print random node ID (by default, random target)
   random-mountpath  print a random mountpath from a given target
```

AIS CLI features a number of miscellaneous and advanced-usage commands.

## Table of Contents
- [Generate shards](#generate-shards)
- [Manually Resilver](#manually-resilver)
- [Preload bucket](#preload-bucket)
- [Remove node from Smap](#remove-node-from-smap)

## Generate shards

`ais advanced gen-shards "BUCKET/TEMPLATE.EXT"`

Put randomly generated shards that can be used for dSort testing.
The `TEMPLATE` must be bash-like brace expansion (see examples) and `.EXT` must be one of: `.tar`, `.tar.gz`.

**Warning**: Remember to always quote the argument (`"..."`) otherwise the brace expansion will happen in terminal.

### Options

| Flag | Type | Description | Default |
| --- | --- | --- | --- |
| `--fsize` | `string` | Single file size inside the shard, can end with size suffix (k, MB, GiB, ...) | `1024`  (`1KB`)|
| `--fcount` | `int` | Number of files inside single shard | `5` |
| `--cleanup` | `bool` | When set, the old bucket will be deleted and created again | `false` |
| `--conc` | `int` | Limits number of concurrent `PUT` requests and number of concurrent shards created | `10` |

### Examples

#### Generate shards with custom number of files and sizes

Generate 10 shards each containing 100 files of size 256KB and put them inside `ais://dsort-testing` bucket (creates it if it does not exist).
Shards will be named: `shard-0.tar`, `shard-1.tar`, ..., `shard-9.tar`.

```console
$ ais advanced gen-shards "ais://dsort-testing/shard-{0..9}.tar" --fsize 262144 --fcount 100
Shards created: 10/10 [==============================================================] 100 %
$ ais ls ais://dsort-testing
NAME		SIZE		VERSION
shard-0.tar	25.05MiB	1
shard-1.tar	25.05MiB	1
shard-2.tar	25.05MiB	1
shard-3.tar	25.05MiB	1
shard-4.tar	25.05MiB	1
shard-5.tar	25.05MiB	1
shard-6.tar	25.05MiB	1
shard-7.tar	25.05MiB	1
shard-8.tar	25.05MiB	1
shard-9.tar	25.05MiB	1
```

#### Generate shards with custom naming template

Generates 100 shards each containing 5 files of size 256KB and put them inside `dsort-testing` bucket.
Shards will be compressed and named: `super_shard_000_last.tgz`, `super_shard_001_last.tgz`, ..., `super_shard_099_last.tgz`

```console
$ ais advanced gen-shards "ais://dsort-testing/super_shard_{000..099}_last.tar" --fsize 262144 --cleanup
Shards created: 100/100 [==============================================================] 100 %
$ ais ls ais://dsort-testing
NAME				SIZE	VERSION
super_shard_000_last.tgz	1.25MiB	1
super_shard_001_last.tgz	1.25MiB	1
super_shard_002_last.tgz	1.25MiB	1
super_shard_003_last.tgz	1.25MiB	1
super_shard_004_last.tgz	1.25MiB	1
super_shard_005_last.tgz	1.25MiB	1
super_shard_006_last.tgz	1.25MiB	1
super_shard_007_last.tgz	1.25MiB	1
...
```

## Manually Resilver

`ais advanced resilver [TARGET_ID]`

Start resilvering objects across all drives on one or all targets.
If `TARGET_ID` is specified, only that node will be resilvered. Otherwise, all targets will be resilvered.

### Examples

```console
$ ais advanced resilver # all targets will be resilvered
Started resilver "NGxmOthtE", use 'ais show job xaction NGxmOthtE' to monitor the progress

$ ais advanced resilver BUQOt8086  # resilver a single node
Started resilver "NGxmOthtE", use 'ais show job xaction NGxmOthtE' to monitor the progress
```

## Preload bucket

`ais advanced preload BUCKET`

Preload bucket's objects metadata into in-memory caches.

### Examples

```console
$ ais advanced preload ais://bucket
```

## Remove node from Smap

`ais advanced remove-from-smap NODE_ID`

Immediately remove node from the cluster map.

Beware! When the node in question is ais target, the operation may (and likely will) result in a data loss that cannot be undone. Use decommission and start/stop maintenance operations to perform graceful removal.

Any attempt to remove from the cluster map `primary` - ais gateway that currently acts as the primary (aka leader) - will fail.

### Examples

```console
$ ais show cluster proxy
PROXY            MEM USED %      MEM AVAIL       UPTIME
BcnQp8083        0.17%           31.12GiB        6m50s
xVMNp8081        0.16%           31.12GiB        6m50s
MvwQp8080[P]     0.18%           31.12GiB        6m40s
NnPLp8082        0.16%           31.12GiB        6m50s


$ ais advanced remove-from-smap MvwQp8080
Node MvwQp 8080 is primary: cannot remove

$ ais advanced remove-from-smap p[xVMNp8081]
$ ais show cluster proxy
PROXY            MEM USED %      MEM AVAIL       UPTIME
BcnQp8083        0.16%           31.12GiB        8m
NnPLp8082        0.16%           31.12GiB        8m
MvwQp8080[P]     0.19%           31.12GiB        7m50s
```

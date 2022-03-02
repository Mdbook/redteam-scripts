# Redteam Scripts
Various Red Team scripts I'm working on.
These scripts are for EDUCATIONAL PURPOSES ONLY. I am not responsible for any damage done due to misuse of these scripts.

## Basic
Basic red team scripts for initial deployment

## Deploy Master
Tool to deploy `basic`, `editor_shim`, `ls_shim`, `passwd_shim`, and `service_herring` to many devices on the network.

## Editor Shim
Shims the `vi`, `vim`, and `nano` binaries, establishing a reverse shell each time they are run.

## Ls Shim
Shims the `ls` binary, establishing a reverse shell each time it is run.

## Passwd Shim
Very simple shim that copies any passwords changed with the `passwd` command to `/tmp/passwords`.

## Service Herring
Installs 32 services with randomized payloads and parameters.

## Shim Handler
Tool to handle the reverse shell connections for `ls_shim` and `editor_shim`.
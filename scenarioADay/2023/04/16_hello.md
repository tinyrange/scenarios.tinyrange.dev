---
title: Welcome to Scenario a Day
description: Writing a Scenario a Day to teach Cyber Security and learn about TinyRange.
date: 16/04/2023
url: /scenarioaday/2023/04/16/hello.html
tags:
  - Programming
  - TinyRange
---

Welcome to my Scenario a Day series where I will be writing a scenario demonstrating some TinyRange or Cyber Security skill each day.

We'll start simple today with a basic demonstration of the scenario lesson framework and a description of the basics of scenario scripts.

Scenarios in TinyRange are written in [Starlark](https://github.com/google/starlark-go). Starlark is a [Python](https://www.python.org/) dialect developed by Google as a scripting language for [Bazel](https://bazel.build/) and other tools. While the language is visually similar to Python there are some important differences.

- Top level variables are read-only (they can't be modified after the script has finished running).
- There are no classes or packages.
- Regular Python code like `import os;os.exit(0);` won't work.

There are two declarations in this file. A virtual machine called `lesson_vm` and the lesson. The lesson starts a web server and provides browser based access to the VM over SSH.

The virtual machine has a custom message of the day (MOTD) that appears when you login to the system and a unprivileged `student` user account.

Both the VM and the lesson are started with the scenario and the lesson opens a web browser automatically when the scenario is booted.

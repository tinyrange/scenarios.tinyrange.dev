---
title: "Running Python inside TinyRange"
description: Simple example of running Python in TinyRange.
date: 19/04/2023
url: /scenarios/2023/04/19/hello_python.html
tags:
  - Python
  - Linux
---

Virtual Machines created with are just regular Linux virtual machines that can have any Linux software installed on them.

This example shows installing Python 3.

The only major departure from previous samples is doubling the RAM allocation to 256MB using `ram = 256`. VMs in TinyRange execute purely from memory by defualt so we allocate additional for storage.
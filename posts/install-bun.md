---
title: "How to Install Bun"
date: 2024-06-01
tags: ["JavaScript", "Tools", "Bun"]
slug: "install-bun"
description: "Quick guide to installing Bun on Windows, macOS, and Linux."

author:
    name: "Leroy"
    email: "hello@leroy.com"
---

## How to Install Bun

Bun is a fast JavaScript runtime. Here are simple install steps.

### macOS / Linux (recommended)

1. Open a terminal.
2. Run the official install script:

```bash
curl -fsSL https://bun.sh/install | bash
```

3. Follow the installer output and restart your terminal.
4. Verify with `bun --version`.

### Windows (WSL or native)

- Use WSL and follow the Linux steps, OR use the official Windows installer when available.

### Troubleshooting

- If `bun` isn't found, ensure `~/.bun/bin` or the installer path is in your `PATH`.

Now you can run scripts with `bun run` and install packages with `bun add`.

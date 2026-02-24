---
title: "How to Set Up SSH Keys"
date: 2024-06-01
tags: ["SSH", "Security", "DevOps"]
slug: "how-to-setup-ssh-keys"
description: "Generate and use SSH keys for passwordless, secure authentication."

author:
    name: "Leroy"
    email: "hello@leroy.com"
---

## How to Set Up SSH Keys

Create key pairs and add your public key to remote services.

### Steps (Linux / macOS / WSL)

1. Generate keys: `ssh-keygen -t ed25519 -C "your_email@example.com"`
2. Add to agent: `eval "$(ssh-agent -s)"` and `ssh-add ~/.ssh/id_ed25519`
3. Copy the public key: `cat ~/.ssh/id_ed25519.pub` and add to GitHub/servers.

### Tips

- Protect private keys with a passphrase.
- Use `ssh-copy-id user@host` for easy server installation.

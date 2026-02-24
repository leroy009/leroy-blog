---
title: "How to Deploy with Docker"
date: 2024-06-01
tags: ["Docker", "Deployment", "DevOps"]
slug: "how-to-deploy-with-docker"
description: "Simple steps to containerize and run an app using Docker."

author:
    name: "Leroy"
    email: "hello@leroy.com"
---

## How to Deploy with Docker

Containerize a simple app and run it locally.

### Steps

1. Create a `Dockerfile` describing your app.
2. Build the image: `docker build -t myapp .`
3. Run the container: `docker run -p 8080:8080 myapp`

### Tips

- Use `.dockerignore` to keep images small.
- Tag images and push to a registry for remote deployment.

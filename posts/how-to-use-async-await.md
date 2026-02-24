---
title: "How to Use async/await (JavaScript)"
date: 2024-06-01
tags: ["JavaScript", "Async", "HowTo"]
slug: "how-to-use-async-await"
description: "A concise explanation of using async/await for asynchronous code in JavaScript."

author:
    name: "Leroy"
    email: "hello@leroy.com"
---

## How to Use async/await

Simplify promise-based code with async functions.

### Example

```js
async function fetchData() {
	try {
		const res = await fetch("/api/data");
		const data = await res.json();
		return data;
	} catch (err) {
		console.error(err);
	}
}
```

### Tips

- Use `try/catch` for errors.
- Avoid unhandled top-level awaits in older environments.

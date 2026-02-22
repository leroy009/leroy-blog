# TODO

- Add documentation
- Add tests
- Optimize the blogs so that they are not read from the file system every time a request is made. Consider using a
- Aside Parser
- Tailwind CSS
- Deploy
- Github Webhooks
- Optimizing: cache?
- Optimizing: goroutines
- Feature: Specify subdirectory and ignore files
- Feature: Paginationcaching mechanism or loading the blogs into memory at startup. This way, the blogs will only be read from the file system once, and subsequent requests can serve the cached data, improving performance.

### Tailwind wathc

```bash
bunx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/style.css --watch
```

# console-to-logger

## What to change

Replace direct `console.log(...)` calls with `logger.info(...)`.

## Do NOT touch

- Test files excluded by `config.yaml`.
- Comments or documentation.

## Examples

```diff
- console.log("started", port)
+ logger.info("started", port)
```

# console-to-logger-astgrep

## What to change

Replace `console.log(...)` matches with `logger.info(...)`.

## Do NOT touch

- Calls on other console methods.
- Test files excluded by `config.yaml`.

## Examples

```diff
- console.log("ready")
+ logger.info("ready")
```

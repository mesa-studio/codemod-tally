# console-to-logger-ast

## What to change

Replace `console.log(...)` call expressions with `logger.info(...)`.

## Do NOT touch

- Calls on other console methods.
- Test files excluded by `config.yaml`.

## Examples

```diff
- console.log(payload)
+ logger.info(payload)
```

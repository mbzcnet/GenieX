## GenieX CLI

Command-line interface for running AI models locally on **Qualcomm** chipsets. Interfaces with the GenieX core runtime and supports two inference backends: **QAIRT** and **llama.cpp**.

### Logging

`GENIEX_LOG` controls log output across the CLI, the C/C++ SDK, and all language bindings (Go, Python, Android):

| Value   | Emits                                    |
|---------|------------------------------------------|
| `none`  | nothing                                  |
| `error` | errors only                              |
| `warn`  | warnings + errors                        |
| `info`  | info + warnings + errors (**default**)   |
| `debug` | debug + info + warnings + errors         |
| `trace` | everything (requires a debug build)      |

```bash
export GENIEX_LOG="debug"          # bash / zsh
$env:GENIEX_LOG="debug"            # PowerShell
```

`NO_COLOR=1` disables ANSI colors.

### Sliding window (qairt) and history trim

The `qairt` backend has a fixed context length (e.g. 4096 tokens). By default `--sliding-window`
is **on**, so overflow evicts the oldest tokens (above a small anchored prefix) instead of
erroring. Pass `--sliding-window=false` to restore the strict error. `llama_cpp` ignores the
flag (it always context-shifts) but honors `--sliding-window-n-keep` for the shift anchor.

Multi-turn chat also trims application history with `--max-history-turns` (default 32; `0` = unlimited):

```bash
geniex infer <model> --max-history-turns 16 --sliding-window-n-keep 8
```

### Model pull

Pull a model non-interactively:

```bash
geniex pull <model>[:<precision>] --model-type <model-type>
```

Pull from a specific model hub:

```bash
geniex pull <model>
geniex pull <model> --model-hub aihub   # options: aihub, hf, localfs
```

Import a model from the local filesystem:

```bash
# hf download <model> --local-dir /path/to/modeldir
geniex pull <model> --model-hub localfs --local-path /path/to/modeldir
```

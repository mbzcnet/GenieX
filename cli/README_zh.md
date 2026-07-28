## GenieX CLI

用于在 **Qualcomm** 芯片上本地运行 AI 模型的命令行工具。对接 GenieX 核心运行时，并支持两种推理后端：**QAIRT** 与 **llama.cpp**。

### 日志

`GENIEX_LOG` 控制 CLI、C/C++ SDK 以及所有语言绑定（Go、Python、Android）的日志输出：

| 取值    | 输出内容                                 |
|---------|------------------------------------------|
| `none`  | 无输出                                   |
| `error` | 仅错误                                   |
| `warn`  | 警告 + 错误                              |
| `info`  | 信息 + 警告 + 错误（**默认**）           |
| `debug` | 调试 + 信息 + 警告 + 错误                |
| `trace` | 全部输出（需要 debug 构建）              |

```bash
export GENIEX_LOG="debug"          # bash / zsh
$env:GENIEX_LOG="debug"            # PowerShell
```

设置 `NO_COLOR=1` 可关闭 ANSI 颜色。

### 滑动窗口（qairt）与历史裁剪

`qairt` 后端的上下文长度是固定的（例如 4096 tokens）。默认 **开启** `--sliding-window`：超限时驱逐最旧 token（保留锚定前缀）而不是报错。若需严格报错，传 `--sliding-window=false`。`llama_cpp` 忽略该开关（始终 context-shift），但会使用 `--sliding-window-n-keep` 作为 shift 锚定长度。

多轮对话还会用 `--max-history-turns`（默认 32；`0` = 不限制）裁剪应用层历史：

```bash
geniex infer <model> --max-history-turns 16 --sliding-window-n-keep 8
```

### 模型拉取

以非交互方式拉取模型：

```bash
geniex pull <model>[:<precision>] --model-type <model-type>
```

从指定的模型仓库拉取：

```bash
geniex pull <model>
geniex pull <model> --model-hub aihub   # 可选：aihub、hf、localfs
```

从本地文件系统导入模型：

```bash
# hf download <model> --local-dir /path/to/modeldir
geniex pull <model> --model-hub localfs --local-path /path/to/modeldir
```

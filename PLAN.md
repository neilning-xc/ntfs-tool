# NTFS 读写挂载 CLI 工具

## Context

macOS 默认以只读方式挂载 NTFS 格式的移动硬盘。开发一个 Go 命令行工具，借助 FUSE + ntfs-3g 实现 NTFS 分区的读写挂载。

## 技术方案

**语言**: Go  
**CLI 框架**: cobra  
**运行时依赖**: macFUSE/FUSE-T + ntfs-3g  

## 子命令

| 命令 | 功能 |
|------|------|
| `ntfs-tool list` | 列出所有 NTFS 分区及其状态 |
| `ntfs-tool mount <device>` | 以读写模式挂载指定 NTFS 分区 |
| `ntfs-tool mount --all` | 挂载所有 NTFS 分区 |
| `ntfs-tool unmount <device>` | 安全卸载指定分区 |
| `ntfs-tool unmount --all` | 卸载所有已挂载的 NTFS 分区 |

## 文件结构

```
ntfs-tool/
├── go.mod
├── go.sum
├── main.go
├── cmd/
│   ├── root.go            # 根命令
│   ├── list.go            # list 子命令
│   ├── mount.go           # mount 子命令
│   └── unmount.go         # unmount 子命令
└── internal/
    ├── disk.go            # NTFS 分区检测（diskutil plist 解析）
    └── shell.go           # Shell 命令执行工具
```

## 核心实现

### 1. NTFS 分区检测 (`internal/disk.go`)
- 调用 `diskutil list -plist` 获取所有磁盘标识符
- 对每个分区调用 `diskutil info -plist <identifier>` 解析 plist
- 检查 `FilesystemType` 或 `FilesystemName` 是否为 `ntfs`
- 返回 `NTFSPartition` 结构体切片（含设备路径、卷标、大小、挂载点）
- 使用 `howett.net/plist` 解析 plist 格式

### 2. 挂载流程 (`cmd/mount.go`)
1. 检查依赖：FUSE 框架 + ntfs-3g 是否已安装
2. 检查 root 权限（ntfs-3g 需要 sudo）
3. 若分区已以只读方式挂载，先用 `diskutil unmount` 卸载
4. 创建挂载点 `/Volumes/<卷标名>`
5. 执行 `ntfs-3g <device> <mountpoint> -o local,allow_other,auto_xattr,noatime`

### 3. 卸载流程 (`cmd/unmount.go`)
- 支持按设备标识符或挂载点卸载
- 使用 `diskutil unmount <mountpoint>` 卸载

### 4. Shell 工具 (`internal/shell.go`)
- 封装 `exec.Command`，返回 stdout/stderr/error
- 提供友好的错误信息

## 依赖检测路径

**FUSE 框架**:
- `/Library/Filesystems/fuse-t.fs` (FUSE-T)
- `/Library/Filesystems/macfuse.fs` (macFUSE)

**ntfs-3g 二进制**:
- `/opt/homebrew/bin/ntfs-3g` (Apple Silicon)
- `/opt/homebrew/sbin/ntfs-3g`
- `/usr/local/bin/ntfs-3g` (Intel)
- `/usr/local/sbin/ntfs-3g`
- 回退: `which ntfs-3g`

> 工具内未检测到依赖时的提示信息：
> - FUSE 未安装 → `brew install --cask macfuse`
> - ntfs-3g 未安装 → `brew install gromgit/fuse/ntfs-3g-mac`

## 运行时依赖安装指南

### 第 1 步：安装 macFUSE

```bash
brew install --cask macfuse
```

安装后需要允许内核扩展：

**Intel Mac**:
1. 系统设置 → 隐私与安全性 → 底部点击"允许" → 重启

**Apple Silicon Mac (M1/M2/M3/M4)**:
1. 系统设置 → 隐私与安全性 → 底部点击"允许"
2. 关机 → 长按电源键进入恢复模式
3. 选项 → 实用工具 → 启动安全性实用工具
4. 选择启动磁盘 → 安全策略 → 降低安全性 → 勾选"允许已验证开发者的内核扩展"
5. 重启后回到 系统设置 → 隐私与安全性 → 允许 macFUSE

### 第 2 步：安装 ntfs-3g（macOS 专用版）

> 注意：官方 `brew install ntfs-3g` 仅支持 Linux，macOS 需要使用社区 fork。

```bash
brew install gromgit/fuse/ntfs-3g-mac
```

### 验证安装

```bash
kextstat | grep fuse        # 检查 macFUSE 内核扩展
ntfs-3g --version           # 检查 ntfs-3g
```

## 验证方式

1. `go build -o ntfs-tool .` 编译成功
2. `./ntfs-tool list` 在无 NTFS 盘时输出 "未发现 NTFS 分区"
3. 插入 NTFS 硬盘后 `sudo ./ntfs-tool mount --all` 完成读写挂载
4. `sudo ./ntfs-tool unmount --all` 安全卸载

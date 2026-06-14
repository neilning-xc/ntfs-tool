# ntfs-tool

macOS NTFS 读写挂载命令行工具。

macOS 默认以只读方式挂载 NTFS 格式的移动硬盘，本工具借助 macFUSE + ntfs-3g 实现读写挂载。

## 安装

### 通过 Homebrew（推荐）

```bash
brew install neilning-xc/tap/ntfs-tool
```

安装后还需要设置 macFUSE 和 ntfs-3g 依赖，参考下方「依赖安装」部分。

### 从源码编译

需要 Go 1.21+：

```bash
git clone https://github.com/neilning-xc/ntfs-tool.git
cd ntfs-tool
CGO_ENABLED=0 go build -o ntfs-tool .
sudo cp ntfs-tool /usr/local/bin/
```

## 依赖安装

### 1. 安装 macFUSE

```bash
brew install --cask macfuse
```

安装后需要允许内核扩展：

**Intel Mac**：
1. 系统设置 → 隐私与安全性 → 底部点击「允许」→ 重启

**Apple Silicon Mac (M1/M2/M3/M4)**：
1. 系统设置 → 隐私与安全性 → 底部点击「允许」
2. 关机 → 长按电源键进入恢复模式
3. 选项 → 实用工具 → 启动安全性实用工具
4. 选择启动磁盘 → 安全策略 → 降低安全性 → 勾选「允许已验证开发者的内核扩展」
5. 重启后回到 系统设置 → 隐私与安全性 → 允许 macFUSE

### 2. 安装 ntfs-3g

> 注意：官方 `brew install ntfs-3g` 仅支持 Linux，macOS 需要使用社区 fork。

```bash
brew install gromgit/fuse/ntfs-3g-mac
```

### 3. 验证依赖

```bash
kextstat | grep fuse    # 检查 macFUSE 内核扩展是否加载
ntfs-3g --version       # 检查 ntfs-3g 是否安装
```

## 使用

### 列出 NTFS 分区

```bash
ntfs-tool list
```

输出示例：

```
发现 1 个 NTFS 分区:

  /dev/disk4s1
    名称: MyDrive
    大小: 500.1 GB
    状态: 已挂载 → /Volumes/MyDrive
```

### 以读写模式挂载

```bash
# 挂载指定分区
sudo ntfs-tool mount disk4s1

# 或使用完整设备路径
sudo ntfs-tool mount /dev/disk4s1

# 挂载所有 NTFS 分区
sudo ntfs-tool mount --all
```

### 卸载

```bash
# 卸载指定分区
sudo ntfs-tool unmount disk4s1

# 卸载所有 NTFS 分区
sudo ntfs-tool unmount --all
```

## 工作原理

1. 通过 `diskutil` 检测系统中所有 NTFS 分区
2. 卸载 macOS 默认的只读挂载
3. 使用 `ntfs-3g` 以读写模式重新挂载到 `/Volumes/<卷标名>`

## 发布新版本

### 前置条件

- 安装 [GoReleaser](https://goreleaser.com)：`brew install goreleaser`
- 安装 [GitHub CLI](https://cli.github.com)：`brew install gh`
- 完成 GitHub 认证：`gh auth login`

### 发布流程

#### 1. 提交所有改动并推送

```bash
git add .
git commit -m "feat: your changes"
git push origin main
```

#### 2. 打标签并推送

```bash
git tag v0.2.0
git push origin v0.2.0
```

推送标签后，GitHub Actions 会自动执行 GoReleaser 构建 darwin amd64/arm64 二进制，并创建 GitHub Release。

#### 3. 获取 SHA256 校验值

等待 GitHub Actions 完成后，下载校验文件：

```bash
gh release download v0.2.0 --pattern "checksums.txt" --output -
```

输出示例：

```
<amd64_sha256>  ntfs-tool_0.2.0_darwin_amd64.tar.gz
<arm64_sha256>  ntfs-tool_0.2.0_darwin_arm64.tar.gz
```

#### 4. 更新 Homebrew Formula

编辑 [homebrew-tap](https://github.com/neilning-xc/homebrew-tap) 仓库中的 `Formula/ntfs-tool.rb`，更新 `version` 和两个 `sha256` 值：

```ruby
version "0.2.0"

# arm64
sha256 "<arm64_sha256>"

# amd64
sha256 "<amd64_sha256>"
```

提交并推送到 `homebrew-tap` 仓库即可。

#### 5. 验证安装

```bash
brew update
brew upgrade ntfs-tool
# 或首次安装
brew install neilning-xc/tap/ntfs-tool
```

## 常见问题

**Q: 提示「未检测到 FUSE 框架」**

安装 macFUSE 并允许内核扩展后重启，参考上方「依赖安装」部分。

**Q: 提示「未找到 ntfs-3g」**

运行 `brew install gromgit/fuse/ntfs-3g-mac` 安装。

**Q: 提示「需要 root 权限」**

挂载和卸载操作需要 sudo：`sudo ntfs-tool mount --all`。

## 卸载

### 卸载 ntfs-tool

```bash
# 通过 Homebrew 安装的
brew uninstall ntfs-tool

# 手动安装的
sudo rm /usr/local/bin/ntfs-tool
```

### 卸载依赖

```bash
# 卸载 ntfs-3g
brew uninstall gromgit/fuse/ntfs-3g-mac

# 卸载 macFUSE
brew uninstall --cask macfuse
```

如果使用的是 FUSE-T：
```bash
brew uninstall --cask fuse-t
```

卸载 macFUSE/FUSE-T 后建议重启，确保内核扩展完全移除。

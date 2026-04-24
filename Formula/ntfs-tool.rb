class NtfsTool < Formula
  desc "macOS NTFS read-write mount CLI tool using macFUSE + ntfs-3g"
  homepage "https://github.com/neilning-xc/ntfs-tool"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/neilning-xc/ntfs-tool/releases/download/v#{version}/ntfs-tool_#{version}_darwin_arm64.tar.gz"
      sha256 "PLACEHOLDER_ARM64_SHA256"
    elsif Hardware::CPU.intel?
      url "https://github.com/neilning-xc/ntfs-tool/releases/download/v#{version}/ntfs-tool_#{version}_darwin_amd64.tar.gz"
      sha256 "PLACEHOLDER_AMD64_SHA256"
    end
  end

  depends_on :macos

  def install
    bin.install "ntfs-tool"
  end

  def caveats
    <<~EOS
      ntfs-tool 需要 macFUSE 和 ntfs-3g 才能正常工作。

      安装依赖:
        brew install --cask macfuse
        brew install gromgit/fuse/ntfs-3g-mac

      安装 macFUSE 后需要在「系统设置 → 隐私与安全性」中允许内核扩展并重启。
      Apple Silicon Mac 还需要进入恢复模式降低安全性设置。
      详见: https://github.com/neilning-xc/ntfs-tool#依赖安装
    EOS
  end

  test do
    assert_match "macOS NTFS", shell_output("#{bin}/ntfs-tool --help")
  end
end

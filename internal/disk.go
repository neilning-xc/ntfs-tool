package internal

import (
	"fmt"
	"os"
	"strings"

	"howett.net/plist"
)

type NTFSPartition struct {
	Identifier string
	DeviceNode string
	VolumeName string
	Size       int64
	MountPoint string
}

func (p NTFSPartition) IsMounted() bool {
	return p.MountPoint != ""
}

func (p NTFSPartition) SizeString() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case p.Size >= TB:
		return fmt.Sprintf("%.1f TB", float64(p.Size)/float64(TB))
	case p.Size >= GB:
		return fmt.Sprintf("%.1f GB", float64(p.Size)/float64(GB))
	case p.Size >= MB:
		return fmt.Sprintf("%.1f MB", float64(p.Size)/float64(MB))
	default:
		return fmt.Sprintf("%d bytes", p.Size)
	}
}

type diskListPlist struct {
	AllDisks   []string `plist:"AllDisks"`
	WholeDisks []string `plist:"WholeDisks"`
}

type diskInfoPlist struct {
	FilesystemType string `plist:"FilesystemType"`
	FilesystemName string `plist:"FilesystemName"`
	VolumeName     string `plist:"VolumeName"`
	DeviceNode     string `plist:"DeviceNode"`
	TotalSize      int64  `plist:"TotalSize"`
	MountPoint     string `plist:"MountPoint"`
}

func ListNTFSPartitions() ([]NTFSPartition, error) {
	result, err := Run("/usr/sbin/diskutil", "list", "-plist")
	if err != nil {
		return nil, fmt.Errorf("无法获取磁盘列表: %w", err)
	}

	var dl diskListPlist
	if _, err := plist.Unmarshal([]byte(result.Stdout), &dl); err != nil {
		return nil, fmt.Errorf("解析磁盘列表失败: %w", err)
	}

	wholeDisks := make(map[string]bool)
	for _, d := range dl.WholeDisks {
		wholeDisks[d] = true
	}

	var partitions []NTFSPartition
	for _, id := range dl.AllDisks {
		if wholeDisks[id] {
			continue
		}
		if p, err := getNTFSInfo(id); err == nil && p != nil {
			partitions = append(partitions, *p)
		}
	}
	return partitions, nil
}

func getNTFSInfo(identifier string) (*NTFSPartition, error) {
	result, err := Run("/usr/sbin/diskutil", "info", "-plist", identifier)
	if err != nil {
		return nil, err
	}

	var info diskInfoPlist
	if _, err := plist.Unmarshal([]byte(result.Stdout), &info); err != nil {
		return nil, err
	}

	fsType := strings.ToLower(info.FilesystemType)
	fsName := strings.ToLower(info.FilesystemName)
	if fsType != "ntfs" && fsName != "ntfs" {
		return nil, nil
	}

	volumeName := info.VolumeName
	if volumeName == "" {
		volumeName = "NTFS"
	}

	return &NTFSPartition{
		Identifier: identifier,
		DeviceNode: info.DeviceNode,
		VolumeName: volumeName,
		Size:       info.TotalSize,
		MountPoint: info.MountPoint,
	}, nil
}

func FindNTFS3G() (string, error) {
	candidates := []string{
		"/opt/homebrew/bin/ntfs-3g",
		"/opt/homebrew/sbin/ntfs-3g",
		"/usr/local/bin/ntfs-3g",
		"/usr/local/sbin/ntfs-3g",
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	// fallback: which
	if result, err := Run("which", "ntfs-3g"); err == nil {
		path := strings.TrimSpace(result.Stdout)
		if path != "" {
			return path, nil
		}
	}
	return "", fmt.Errorf("未找到 ntfs-3g，请先安装: brew install gromgit/fuse/ntfs-3g-mac")
}

func CheckFUSE() error {
	fusePaths := []string{
		"/Library/Filesystems/macfuse.fs",
		"/Library/Filesystems/fuse-t.fs",
	}
	for _, path := range fusePaths {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
	}
	return fmt.Errorf("未检测到 FUSE 框架，请先安装: brew install --cask macfuse")
}

func IsRoot() bool {
	return os.Geteuid() == 0
}

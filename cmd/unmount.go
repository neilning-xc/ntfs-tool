package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/neilning-xc/ntfs-tool/internal"
	"github.com/spf13/cobra"
)

var unmountAll bool

var unmountCmd = &cobra.Command{
	Use:     "unmount [device]",
	Aliases: []string{"umount"},
	Short:   "卸载 NTFS 分区",
	Long: `安全卸载已挂载的 NTFS 分区。

示例:
  sudo ntfs-tool unmount disk4s1
  sudo ntfs-tool unmount /dev/disk4s1
  sudo ntfs-tool unmount --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		partitions, err := internal.ListNTFSPartitions()
		if err != nil {
			return err
		}

		mounted := filterMounted(partitions)
		if len(mounted) == 0 {
			fmt.Println("没有已挂载的 NTFS 分区")
			return nil
		}

		if unmountAll {
			for _, p := range mounted {
				if err := unmountPartition(p); err != nil {
					fmt.Fprintf(os.Stderr, "卸载 %s 失败: %v\n", p.DeviceNode, err)
				}
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("请指定设备标识符，或使用 --all 卸载所有 NTFS 分区")
		}

		device := args[0]
		if !strings.HasPrefix(device, "/dev/") {
			device = "/dev/" + device
		}

		var target *internal.NTFSPartition
		for _, p := range mounted {
			if p.DeviceNode == device {
				target = &p
				break
			}
		}
		if target == nil {
			return fmt.Errorf("未找到已挂载的 NTFS 分区: %s", args[0])
		}
		return unmountPartition(*target)
	},
}

func unmountPartition(p internal.NTFSPartition) error {
	fmt.Printf("正在卸载: %s (%s) ...\n", p.MountPoint, p.VolumeName)
	if _, err := internal.Run("/usr/sbin/diskutil", "unmount", p.MountPoint); err != nil {
		return fmt.Errorf("卸载失败: %w", err)
	}
	fmt.Printf("已卸载: %s\n", p.VolumeName)
	return nil
}

func filterMounted(partitions []internal.NTFSPartition) []internal.NTFSPartition {
	var result []internal.NTFSPartition
	for _, p := range partitions {
		if p.IsMounted() {
			result = append(result, p)
		}
	}
	return result
}

func init() {
	unmountCmd.Flags().BoolVarP(&unmountAll, "all", "a", false, "卸载所有 NTFS 分区")
}

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/neilning-xc/ntfs-tool/internal"
	"github.com/spf13/cobra"
)

var mountAll bool

var mountCmd = &cobra.Command{
	Use:   "mount [device]",
	Short: "以读写模式挂载 NTFS 分区",
	Long: `以读写模式挂载 NTFS 分区。需要 sudo 权限。

示例:
  sudo ntfs-tool mount disk4s1
  sudo ntfs-tool mount /dev/disk4s1
  sudo ntfs-tool mount --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.CheckFUSE(); err != nil {
			return err
		}
		if !internal.IsRoot() {
			return fmt.Errorf("挂载 NTFS 分区需要 root 权限，请使用: sudo ntfs-tool mount ...")
		}

		ntfs3gPath, err := internal.FindNTFS3G()
		if err != nil {
			return err
		}

		partitions, err := internal.ListNTFSPartitions()
		if err != nil {
			return err
		}
		if len(partitions) == 0 {
			fmt.Println("未发现 NTFS 分区")
			return nil
		}

		if mountAll {
			for _, p := range partitions {
				if err := mountPartition(p, ntfs3gPath); err != nil {
					fmt.Fprintf(os.Stderr, "挂载 %s 失败: %v\n", p.DeviceNode, err)
				}
			}
			return nil
		}

		if len(args) == 0 {
			return fmt.Errorf("请指定设备标识符，或使用 --all 挂载所有 NTFS 分区")
		}

		device := args[0]
		if !strings.HasPrefix(device, "/dev/") {
			device = "/dev/" + device
		}

		var target *internal.NTFSPartition
		for _, p := range partitions {
			if p.DeviceNode == device {
				target = &p
				break
			}
		}
		if target == nil {
			return fmt.Errorf("未找到 NTFS 分区: %s", args[0])
		}
		return mountPartition(*target, ntfs3gPath)
	},
}

func mountPartition(p internal.NTFSPartition, ntfs3gPath string) error {
	mountPoint := "/Volumes/" + p.VolumeName

	if p.IsMounted() {
		fmt.Printf("正在卸载只读挂载: %s ...\n", p.DeviceNode)
		if _, err := internal.Run("/usr/sbin/diskutil", "unmount", p.DeviceNode); err != nil {
			return fmt.Errorf("卸载失败: %w", err)
		}
	}

	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("创建挂载点失败: %w", err)
	}

	fmt.Printf("正在以读写模式挂载: %s → %s ...\n", p.DeviceNode, mountPoint)
	options := fmt.Sprintf("local,allow_other,auto_xattr,noatime,volname=%s", p.VolumeName)
	if _, err := internal.Run(ntfs3gPath, p.DeviceNode, mountPoint, "-o", options); err != nil {
		return fmt.Errorf("ntfs-3g 挂载失败: %w", err)
	}

	fmt.Printf("挂载成功: %s (%s)\n", p.VolumeName, p.SizeString())
	return nil
}

func init() {
	mountCmd.Flags().BoolVarP(&mountAll, "all", "a", false, "挂载所有 NTFS 分区")
}

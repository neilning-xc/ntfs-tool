package cmd

import (
	"fmt"

	"github.com/neilning-xc/ntfs-tool/internal"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有 NTFS 分区",
	RunE: func(cmd *cobra.Command, args []string) error {
		partitions, err := internal.ListNTFSPartitions()
		if err != nil {
			return err
		}
		if len(partitions) == 0 {
			fmt.Println("未发现 NTFS 分区")
			return nil
		}

		fmt.Printf("发现 %d 个 NTFS 分区:\n\n", len(partitions))
		for _, p := range partitions {
			status := "未挂载"
			if p.IsMounted() {
				status = fmt.Sprintf("已挂载 → %s", p.MountPoint)
			}
			fmt.Printf("  %s\n", p.DeviceNode)
			fmt.Printf("    名称: %s\n", p.VolumeName)
			fmt.Printf("    大小: %s\n", p.SizeString())
			fmt.Printf("    状态: %s\n\n", status)
		}
		return nil
	},
}

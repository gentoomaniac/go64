package cmd

import (
	"fmt"

	"github.com/gentoomaniac/go64/c64"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the emulator",
	Long:  `long description`,
	Run: func(cmd *cobra.Command, args []string) {
		system := &c64.C64{}

		system.Init("rom/basic.rom", "rom/kernal.rom", "rom/character.rom")

		fmt.Println(system.DumpMemory(0, c64.MaxMemoryAddress))
		fmt.Println(system.Mpu.DumpRegisters())
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

package cmd

import (
	"os/user"

	"github.com/blixenkrone/dotgo/linker"
	"github.com/blixenkrone/dotgo/linker/symbolic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var selectDotFilesCmd = &cobra.Command{
	Use:   "ln",
	Short: "Create a symbolic link for files",
	// Args:  cobra.MinimumNArgs(1),
	Run: linkFromDir(),
}

var (
	sourceDirFlag, destDirFlag string
	recursiveFlag              bool
	whitelistedFlag            []string
)

func linkFromDir() CobraFunc {
	l := logrus.New()
	return func(cmd *cobra.Command, _ []string) {
		var srcdir, destdir string
		if sourceDirFlag == "" {
			cmd.PrintErr("no source directory specified")
			return
		}
		srcdir = sourceDirFlag

		if destDirFlag != "" {
			destdir = destDirFlag
		} else {
			usr, err := user.Current()
			if err != nil {
				cmd.PrintErrf("user was not found: %v", err)
				return
			}
			destdir = usr.HomeDir
		}

		o, err := symbolic.NewDirOp(l, srcdir, destdir, recursiveFlag)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		unlinked, total, err := linker.Link(o)
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		cmd.Printf("linked %d/%d files to %s", unlinked, total, destdir)
	}
}

func initSelectDotfilesCmd() {
	rootCmd.AddCommand(selectDotFilesCmd)
	selectDotFilesCmd.PersistentFlags().StringVarP(&sourceDirFlag, "source", "s", "", "Source directory to read from")
	selectDotFilesCmd.PersistentFlags().StringVarP(&destDirFlag, "dest", "d", "", "Destination directory to write the link to")
	selectDotFilesCmd.PersistentFlags().BoolVarP(&recursiveFlag, "recursive", "r", true, "If the linking should be done on nested folders in source dir provided")
	selectDotFilesCmd.PersistentFlags().StringArrayVarP(&whitelistedFlag, "whitelist", "w", nil, "Files to link regardless of them being a dotfile - duplicates are ignored")
}

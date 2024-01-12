/*
Copyright © 2021 Meir Gabay <unfor19@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bytes"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/mikefarah/yq/v4/cmd"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Reads a YAML file and parses the YAML anchors as plain text in a new file",
	Long: `Reads a YAML file and parses the YAML anchors as plain text in a new file
For example:

yarser parse --watch "$YARSER_SRC_FILE_PATH" "$YARSER_DST_FILE_PATH"
yarser parse "$YARSER_SRC_FILE_PATH"
	
	`,
	// srcFilePath, dstFilePath
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var srcFilePath string = args[0]
		var dstFilePath string = args[1]
		watch, _ := cmd.Flags().GetBool("watch")
		if watch {
			CustomWatcher(srcFilePath, dstFilePath, parseYaml)
		} else {
			parseYaml(srcFilePath, dstFilePath)
		}
	},
}

func getRootCommand() *cobra.Command {
	return cmd.New()
}

func runYq(input string) string {
	/*
	 * TODO - use use the built in https://pkg.go.dev/github.com/mikefarah/yq/v4@v4.40.5/pkg/yqlib
	 * and don't call out to the yq	 executable if possible.
	 */
	cmd := getRootCommand()
	buffer := new(bytes.Buffer)

	cmd.SetOutput(buffer)
	cmd.SetArgs(strings.Split(input, " "))

	err := cmd.Execute()
	if err != nil {
		logger.Println(err.Error())
		return ""
	}

	logger.Debug("Successfully executed:", "yq", input)
	output := buffer.String()

	return stripansi.Strip(output)
}

func parseYaml(srcFilePath string, dstFilePath string) error {
	OutputFile := CreateEmptyFile(dstFilePath)
	defer OutputFile.Close()
	logger.Debug("Destination file path", dstFilePath)
	explodeResult := runYq("eval-all" + " explode(.) " + srcFilePath)
	OutputFile.Write([]byte(explodeResult))

	delResult := runYq("eval-all" + " --inplace" + " del(.\".*\") " + dstFilePath)
	if delResult == "" {
		logger.Debug("Deleted all nodes that start with '.'")
	} else {
		logger.Debug("Deletion result: ", delResult)
	}

	logger.Info("Successfully parsed ", srcFilePath, " to ", dstFilePath)

	return nil
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	parseCmd.Flags().BoolP("watch", "w", false, "Watch file for changes")
}

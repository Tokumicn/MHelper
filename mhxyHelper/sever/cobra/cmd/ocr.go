package cmd

import (
	"github.com/spf13/cobra"
	"mhxyHelper/internal/service/buildQA"
)

func init() {

}

var ocrCmd = &cobra.Command{
	Use:   "ocr",
	Short: "OCR构建商品信息",
	Long:  "OCR构建商品信息",
	Run: func(cmd *cobra.Command, args []string) {
		buildQA.BuildFromOCR()
	},
}

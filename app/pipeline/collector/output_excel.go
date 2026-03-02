package collector

import (
	"fmt"
	"os"

	"github.com/andeya/pholcus/common/util"
	"github.com/andeya/pholcus/common/xlsx"
	"github.com/andeya/pholcus/config"
	"github.com/andeya/pholcus/logs"
	"github.com/andeya/pholcus/runtime/cache"
)

// --- Excel Output ---

func init() {
	DataOutput["excel"] = func(self *Collector) (err error) {
		defer func() {
			if p := recover(); p != nil {
				err = fmt.Errorf("%v", p)
			}
		}()

		var (
			file   *xlsx.File
			row    *xlsx.Row
			cell   *xlsx.Cell
			sheets = make(map[string]*xlsx.Sheet)
		)

		file = xlsx.NewFile()

		for _, datacell := range self.dataDocker {
			var subNamespace = util.FileNameReplace(self.subNamespace(datacell))
			if _, ok := sheets[subNamespace]; !ok {
				sheet, err := file.AddSheet(subNamespace)
				if err != nil {
					logs.Log.Error("%v", err)
					continue
				}
				sheets[subNamespace] = sheet
				row = sheets[subNamespace].AddRow()
				for _, title := range self.MustGetRule(datacell["RuleName"].(string)).ItemFields {
					row.AddCell().Value = title
				}
				if self.Spider.OutDefaultField() {
					row.AddCell().Value = "当前链接"
					row.AddCell().Value = "上级链接"
					row.AddCell().Value = "下载时间"
				}
			}

			row = sheets[subNamespace].AddRow()
			for _, title := range self.MustGetRule(datacell["RuleName"].(string)).ItemFields {
				cell = row.AddCell()
				vd := datacell["Data"].(map[string]interface{})
				if v, ok := vd[title].(string); ok || vd[title] == nil {
					cell.Value = v
				} else {
					cell.Value = util.JsonString(vd[title])
				}
			}
			if self.Spider.OutDefaultField() {
				row.AddCell().Value = datacell["Url"].(string)
				row.AddCell().Value = datacell["ParentUrl"].(string)
				row.AddCell().Value = datacell["DownloadTime"].(string)
			}
		}
		folder := config.TEXT_DIR + "/" + cache.StartTime.Format("2006-01-02 150405")
		filename := fmt.Sprintf("%v/%v__%v-%v.xlsx", folder, util.FileNameReplace(self.namespace()), self.sum[0], self.sum[1])

		// Create directory if needed
		f2, err := os.Stat(folder)
		if err != nil || !f2.IsDir() {
			if err := os.MkdirAll(folder, 0777); err != nil {
				logs.Log.Error("Error: %v\n", err)
			}
		}

		err = file.Save(filename)
		return
	}
}

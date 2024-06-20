package excel

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/stretchr/testify/require"
	"github.com/tealeg/xlsx"
)

func getTealegXlsx(t *testing.T) *xlsx.File {
	titles := []string{
		"name",
		"age",
		"score",
		"weight",
	}
	type row struct {
		Name   string
		Age    int
		Score  int
		Weight float64
	}
	datas := []*row{
		{
			Name:   "tsukiyo",
			Age:    18,
			Score:  100,
			Weight: 126,
		},
		{
			Name:   "lazywoo",
			Age:    23,
			Score:  61,
			Weight: 129,
		},
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	require.NoError(t, err)
	titleRow := sheet.AddRow()
	for _, title := range titles {
		cell := titleRow.AddCell()
		cell.Value = title
		cell.GetStyle().Font.Color = "00FF0000"
	}
	// 插入内容
	for _, v := range datas {
		row := sheet.AddRow()
		//if slice, ok := v.([]string); ok {
		//	row.WriteSlice(&slice, -1)
		//} else {
		//	row.WriteStruct(v, -1)
		//}
		row.WriteStruct(v, -1)
	}

	return file
}

func genFileName(name string) string {
	return url.QueryEscape(time.Now().Format(time.DateTime) + "工作需求" + ".xlsx")
}

// use tealeg/xlsx
func TestGenExcelV1(t *testing.T) {
	file := getTealegXlsx(t)
	fileName := genFileName("工作需求")
	err := file.Save(fileName)
	require.NoError(t, err)
}

// 这里这种io.Pipe没有意义，但是没有坏处，并且有拓展性
func TestHTTPExportUsage(t *testing.T) {
	engine := gin.Default()
	engine.GET("/export", func(c *gin.Context) {
		pr, pw := io.Pipe()

		var wg sync.WaitGroup
		wg.Add(1)

		// sender
		go func() {
			defer wg.Done()

			if _, err := io.Copy(c.Writer, pr); err != nil {
				t.Error("failed to close io copy")
			}

			if err := pr.Close(); err != nil {
				t.Error("pipe reader close failed")
			}
		}()

		// generator
		func() {
			c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", genFileName("工作需求")))

			file := getTealegXlsx(t)
			err := file.Write(pw)
			if err != nil {
				t.Error("excel file error")
			}
		}()

		wg.Wait()

		c.Status(http.StatusOK)
	})
	require.NoError(t, engine.Run(":8080"))
}

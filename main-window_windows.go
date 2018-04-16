package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lxn/walk"
	decl "github.com/lxn/walk/declarative"
)

func showMainWindow(args []string) {
	var err error
	var inTE *walk.TextEdit
	var inFiles []string = []string{}
	var logView *LogView
	var logViewParent *walk.Splitter
	var mainWin *walk.MainWindow
	var startOptimizeBtn *walk.PushButton
	var cancelOptimizeBtn *walk.PushButton
	var doCancel context.CancelFunc

	ensureLibExeExists(102, cfg.Jpegoptim)
	ensureLibExeExists(103, cfg.Pngquant)

	isCanceled := false
	isProcessing := false

	onClickStartOptimizeBtn := func() {
		if isProcessing {
			walk.MsgBox(mainWin, "提示", "正在进行优化，请稍等片刻。", walk.MsgBoxIconWarning)
			return
		}

		inFiles = uniqueLines(inFiles)
		inTE.SetText(strings.Join(inFiles, "\r\n"))

		go func() {
			defer func() {
				r := recover()

				startOptimizeBtn.SetText("开始优化")
				startOptimizeBtn.SetEnabled(true)
				cancelOptimizeBtn.SetVisible(false)
				isProcessing = false

				if r != nil {
					if isCanceled {
						log.Print("已经取消。")
					} else {
						log.Print("优化时遇到错误：" + formatError(r))
					}
				}

			}()

			var ctx context.Context
			ctx, doCancel = context.WithCancel(context.Background())

			isProcessing = true
			isCanceled = false
			cancelOptimizeBtn.SetVisible(true)
			cancelOptimizeBtn.SetEnabled(true)
			startOptimizeBtn.SetEnabled(false)
			startOptimizeBtn.SetText("正在优化...")
			log.Print("开始优化...")

			for i := 0; i < len(inFiles); i++ {
				x := inFiles[i]
				doOptimize(x, ctx)
			}

			log.Print("优化完成。")
		}()
	}

	onClickCancelOptimizeBtn := func() {
		if doCancel != nil {
			isCanceled = true
			doCancel()
			doCancel = nil
			cancelOptimizeBtn.SetEnabled(false)
		}
	}

	mainWinDef := decl.MainWindow{
		Title:    "图片优化工具",
		AssignTo: &mainWin,
		MinSize:  decl.Size{600, 400},
		Layout:   decl.VBox{},
		OnDropFiles: func(files []string) {
			inFiles = append(inFiles, files...)
			inTE.SetText(strings.Join(inFiles, "\r\n"))
			log.Printf("已经收到文件，请点击【开始优化】来开始优化。")
		},
		Children: []decl.Widget{
			decl.VSplitter{
				HandleWidth: 3,
				// Background:  decl.SolidColorBrush{Color: walk.RGB(255, 0, 0)},
				Children: []decl.Widget{
					decl.Composite{
						Layout: decl.Grid{Columns: 10},
						Children: []decl.Widget{
							decl.Label{
								ColumnSpan: 10,
								Text:       "要处理的图片文件/文件夹/压缩包列表：(请从 Windows 资源管理器往这里拖放文件...)",
							},
							decl.TextEdit{
								ColumnSpan:    10,
								AssignTo:      &inTE,
								ReadOnly:      false,
								Text:          "",
								StretchFactor: 1,
								OnTextChanged: func() {
									text := inTE.Text()
									inFiles = filterEmptyLines(splitIntoLines(text))
								},
							},
							decl.PushButton{
								ColumnSpan: 3,
								AssignTo:   &startOptimizeBtn,
								Text:       "开始优化",
								MaxSize:    decl.Size{Width: 600, Height: 40},
								OnClicked:  onClickStartOptimizeBtn,
							},
							decl.VSpacer{
								ColumnSpan: 4,
							},
							decl.PushButton{
								ColumnSpan: 3,
								AssignTo:   &cancelOptimizeBtn,
								Text:       "取消",
								MaxSize:    decl.Size{Width: 600, Height: 40},
								Visible:    false,
								OnClicked:  onClickCancelOptimizeBtn,
							},
						},
					},
					decl.Composite{
						Layout: decl.Grid{Columns: 10},
						Children: []decl.Widget{
							decl.Label{
								ColumnSpan: 10,
								Text:       "优化日志：",
							},
							decl.VSplitter{
								ColumnSpan:    10,
								StretchFactor: 1,
								// Background:    decl.SolidColorBrush{Color: walk.RGB(255, 0, 0)},
								AssignTo: &logViewParent,
							},
						},
					},
				},
			},
		},
	}

	if err = mainWinDef.Create(); err != nil {
		log.Fatal(err)
		return
	}

	if logView, err = NewLogView(logViewParent); err != nil {
		log.Fatal(err)
		return
	}

	// GUI application dont need a console
	// C.freeConsole()
	log.SetOutput(newLogWriterOfLogViewAndStdout(logView))

	// set the main window icon
	// 2: defined in app.rc -- the app icon
	if mainIco, err := walk.NewIconFromResourceId(2); err == nil {
		mainWin.SetIcon(mainIco)
	}

	if len(args) > 0 {
		inFiles = uniqueLines(args)
		inTE.SetText(strings.Join(inFiles, "\r\n"))
	}

	mainWin.Run()
}

type tLogWriterOfLogViewAndStdout struct {
	logView *LogView
}

func newLogWriterOfLogViewAndStdout(logView *LogView) io.Writer {
	return &tLogWriterOfLogViewAndStdout{
		logView: logView,
	}
}

func (t *tLogWriterOfLogViewAndStdout) Write(p []byte) (n int, err error) {
	// n, err = os.Stderr.Write(p)
	// if err != nil {
	// 	return
	// }

	n = len(p)
	t.logView.AppendText(string(p[:n]))
	return
}

func splitIntoLines(text string) []string {
	return strings.Split(text, "\n")
}

func filterEmptyLines(lines []string) []string {
	r := make([]string, len(lines))
	n := 0
	for _, x := range lines {
		x = strings.TrimSpace(x)
		if x != "" {
			r[n] = x
			n++
		}
	}

	return r[:n]
}

func uniqueLines(lines []string) []string {
	m := make(map[string]bool)
	r := make([]string, len(lines))
	n := 0
	for _, x := range lines {
		x = strings.TrimSpace(x)
		if x != "" && !m[x] {
			m[x] = true
			r[n] = x
			n++
		}
	}

	return r[:n]
}

func debugDump() {
	walk.MsgBox(nil, "Debug", fmt.Sprintf("argv: %#v", os.Args), walk.MsgBoxOK)
}

func ensureLibExeExists(libResId int, libExePath string) {
	if fileExists(libExePath) {
		return
	}

	extList := strings.Split(os.Getenv("PATHEXT"), ";")
	for _, ext := range extList {
		if fileExists(libExePath + ext) {
			return
		}
	}

	outputLibExePath := libExePath
	if filepath.Ext(libExePath) == "" {
		outputLibExePath += ".exe"
	}

	extractAppResource(libResId, outputLibExePath, 0777)
}

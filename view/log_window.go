package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
	"image/color"
	"strings"
	"we2log/config"
	"we2log/service"
)

// 窗口信息map
var winInfoMap = make(map[string]*winInfo, service.GroupNum)

type winInfo struct {
	// 日志状态是否继续接收
	logStatus bool
	// 日志是否置底
	logBottom bool
	// 日志列表
	list *fyne.Container
}

// logView 创建日志视图
func logView(group string) fyne.CanvasObject {
	// 顶部控件
	// 暂停继续
	pause := widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {})
	pause.OnTapped = func() { logPauseOrContinue(group, pause) }
	// 置底
	bottom := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {})
	bottom.OnTapped = func() { logBottom(group, bottom) }
	hBox := container.NewHBox(pause, bottom)

	// 日志列表
	list := container.NewVBox()
	scroll := container.NewScroll(list)
	scroll.Resize(fyne.NewSize(1024, 607))
	scroll.Move(fyne.NewPos(0, 41))

	// 运行日志监听
	go onChange(group, list, func() {
		list.Refresh()
		if winInfoMap[group].logBottom {
			scroll.ScrollToBottom()
		}
	})

	// 分割线
	separator := widget.NewSeparator()
	separator.Resize(fyne.NewSize(1024, 1))
	separator.Move(fyne.NewPos(0, 40))
	return container.NewWithoutLayout(scroll, separator, hBox)
}

// LogThemeChange 日志改变主题色
func LogThemeChange() {
	// 日志颜色更换
	for _, info := range winInfoMap {
		for _, o := range info.list.Objects {
			text := o.(*canvas.Text)
			if CurrentTheme {
				if text.Color == color.White {
					text.Color = color.Black
				}
			} else {
				if text.Color == color.Black {
					text.Color = color.White
				}
			}
		}
	}
}

// 日志置底
func logBottom(group string, button *widget.Button) {
	info := winInfoMap[group]
	info.logBottom = !info.logBottom
	if info.logBottom {
		button.SetIcon(theme.ContentRemoveIcon())
	} else {
		button.SetIcon(theme.DownloadIcon())
	}
}

// 日志暂停或者继续
func logPauseOrContinue(group string, button *widget.Button) {
	info := winInfoMap[group]
	info.logStatus = !info.logStatus
	if info.logStatus {
		button.SetIcon(theme.MediaPauseIcon())
	} else {
		button.SetIcon(theme.MediaPlayIcon())
	}
}

// 监听日志改变
func onChange(group string, list *fyne.Container, after func()) {
	// 坐标
	point := 0
	// 缓冲行
	lines := config.Yml.Log.Lines
	winInfoMap[group] = &winInfo{
		logStatus: true,
		logBottom: true,
		list:      list,
	}
	for m := range service.MsgChanMap[group] {
		// 是否开启日志
		if !winInfoMap[group].logStatus {
			continue
		}
		text := getText(group, &m)
		if point < lines {
			list.Add(text)
			point++
		} else {
			for i := 0; i < lines; i++ {
				// 如果最后一个元素则改为新数据
				if i == lines-1 {
					list.Objects[i] = text
				} else {
					list.Objects[i] = list.Objects[i+1]
				}
			}
		}
		// 最后执行
		after()
	}
}

// 获取文本及对应颜色
func getText(group string, msg *string) *canvas.Text {
	var c color.Color
	switch {
	case strings.Contains(*msg, "INFO"):
		c = colornames.Green
	case strings.Contains(*msg, "WARN"):
		c = colornames.Yellow
	case strings.Contains(*msg, "ERROR"):
		c = colornames.Red
	default:
		if CurrentTheme {
			c = color.Black
		} else {
			c = color.White
		}
	}

	text := canvas.NewText(*msg, c)
	text.TextSize = config.Yml.Log.FontSize
	return text
}

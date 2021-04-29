package view

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	// 日志窗口是否全屏
	logFull bool
	// 日志列表
	list *fyne.Container
}

// logView 创建日志视图
func logView(group string, w fyne.Window) fyne.CanvasObject {
	// 日志列表
	list := container.NewVBox()
	scroll := container.NewScroll(list)
	scroll.Resize(fyne.NewSize(1024, 648-8))

	// 运行日志监听
	go onChange(group, list, func() {
		if winInfoMap[group].logStatus {
			scroll.ScrollToBottom()
			list.Refresh()
		}
	})

	// 设置窗口事件监听
	w.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		onTypedKey(group, w, event)
	})
	return container.NewMax(scroll)
}

// 按键事件
func onTypedKey(group string, w fyne.Window, event *fyne.KeyEvent) {
	if event.Name == fyne.KeyF {
		logFull(group, w)
	}
	if event.Name == fyne.KeySpace {
		logPauseOrContinue(group, w)
	}
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

// 日志全屏
func logFull(group string, w fyne.Window) {
	info := winInfoMap[group]
	info.logFull = !info.logFull
	w.SetFullScreen(info.logFull)
}

// 日志暂停或者继续
func logPauseOrContinue(group string, w fyne.Window) {
	info := winInfoMap[group]
	info.logStatus = !info.logStatus
	if info.logStatus {
		w.SetTitle(group)
	} else {
		w.SetTitle(group + "(已暂停)")
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
		list:      list,
	}
	for m := range service.MsgChanMap[group] {
		// 是否开启日志
		if !winInfoMap[group].logStatus {
			continue
		}
		text := getText(&m)
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
func getText(msg *string) *canvas.Text {
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

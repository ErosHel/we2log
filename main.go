package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"os"
	"runtime"
	"strings"
	"we2log/config"
	"we2log/model/log"
	"we2log/view"
)

const fontEnv = "FYNE_FONT"

func main() {
	// 初始化配置
	config.InitYaml()
	// 设置字体
	setFont()
	// 建立界面
	a := app.New()
	// 设置主题
	themeChangeListener(a)
	// 构建主窗口
	buildMainWindow(&a)
	// 解除字体
	_ = os.Unsetenv(fontEnv)
}

// 主题切换监听
func themeChangeListener(a fyne.App) {
	view.CurrentTheme = a.Settings().ThemeVariant() == 1
	ch := make(chan fyne.Settings)
	go func() {
		for set := range ch {
			themeSite := fmt.Sprintf("%v", set.Theme())
			if view.CurrentTheme != strings.HasPrefix(themeSite, "&{1") {
				view.CurrentTheme = !view.CurrentTheme
				view.LogThemeChange()
			}
		}
	}()
	a.Settings().AddChangeListener(ch)
}

// 构建主窗口
func buildMainWindow(a *fyne.App) {
	w := (*a).NewWindow("we2log")
	w.Resize(fyne.NewSize(320, 490))
	w.SetContent(view.CreateMainWindow(w, a))
	// 保存配置
	w.SetOnClosed(config.SaveLocalCache)
	w.CenterOnScreen()
	w.SetFixedSize(true)
	w.ShowAndRun()
}

//设置中文字体
func setFont() {
	sysType := runtime.GOOS
	var err error
	if sysType == "darwin" {
		err = os.Setenv(fontEnv, "/System/Library/Fonts/Supplemental/Arial Unicode.ttf")
	}
	if sysType == "windows" {
		err = os.Setenv(fontEnv, "C:/Windows/fonts/simhei.ttf")
	}
	if err != nil {
		log.Error(fmt.Sprintf("字体变量设置错误: %v", err))
	}
}

package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"strings"
	"we2log/config"
	"we2log/resource/icon"
	"we2log/service"
)

var (
	runStatus bool
	// CurrentTheme 当前主题 true白色 false黑色
	CurrentTheme bool
	logWin       []fyne.Window
)

// CreateMainWindow 创建主窗口
func CreateMainWindow(w fyne.Window, a *fyne.App) fyne.CanvasObject {
	// 运行按钮
	runBtn := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
	runBtn.OnTapped = func() { runOrStop(runBtn, a, w) }
	runBtn.Resize(fyne.NewSize(40, 40))
	// 增加分组输入框
	groupEntry := widget.NewEntry()
	groupEntry.SetPlaceHolder("分组名称")
	groupEntry.Resize(fyne.NewSize(80, 40))
	groupEntry.Move(fyne.NewPos(43, 0))
	// 增加分组
	groupAdd := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	groupAdd.Resize(fyne.NewSize(40, 40))
	groupAdd.Move(fyne.NewPos(43+83, 0))
	// 主题切换
	themeBtn := widget.NewButtonWithIcon("", getThemeIcon(), func() {})
	themeBtn.OnTapped = func() { themeChange(themeBtn) }
	themeBtn.Resize(fyne.NewSize(40, 40))
	themeBtn.Move(fyne.NewPos(313-40, 0))
	// 顶部横排控件
	topTitle := container.NewWithoutLayout(runBtn, groupEntry, groupAdd, themeBtn)
	topTitle.Resize(fyne.NewSize(313, 40))

	// 分割线
	separator := widget.NewSeparator()
	separator.Resize(fyne.NewSize(313, 1))
	separator.Move(fyne.NewPos(0, 42))

	list := getGroupBindList()
	groupList := widget.NewListWithData(list, func() fyne.CanvasObject {
		return createGroupItem(w, list)
	}, updateGroupItem)
	groupList.Move(fyne.NewPos(0, 42))
	groupList.Resize(fyne.NewSize(313, 445))
	groupAdd.OnTapped = func() {
		addGroup(groupEntry.Text, w, list)
	}

	return container.NewWithoutLayout(topTitle, separator, groupList)
}

// 添加分组
func addGroup(name string, w fyne.Window, list binding.StringList) {
	if name == "" {
		return
	}
	groups := *config.Yml.Log.Group
	// 判断分组重复
	for _, group := range groups {
		if group.Name == name {
			dialogInfo("分组重复", "已经存在的分组名称", w)
			return
		}
	}
	*config.Yml.Log.Group = append(groups, config.Group{
		Name:  name,
		OnOff: false,
		Ssh:   &[]config.Ssh{},
	})
	_ = list.Set(getGroupNameInfo())
	config.SaveLocalCache()
}

// 提示框
func dialogInfo(title string, content string, w fyne.Window) {
	dl := dialog.NewInformation(title, content, w)
	dl.SetDismissText("返回")
	dl.Show()
}

// 获取分组绑定数据
func getGroupBindList() binding.StringList {
	list := binding.NewStringList()
	_ = list.Set(getGroupNameInfo())
	return list
}

// 获取分组最新信息
func getGroupNameInfo() []string {
	if len(*config.Yml.Log.Group) == 0 {
		return []string{}
	}
	groups := *config.Yml.Log.Group
	ss := make([]string, len(groups))
	for i, group := range groups {
		// [ssh名称][是否选中][分组坐标]
		ss[i] = fmt.Sprintf("%s,%t,%d", group.Name, group.OnOff, i)
	}
	return ss
}

// 创建分组的元素
func createGroupItem(w fyne.Window, list binding.StringList) fyne.CanvasObject {
	// 分组名称
	groupName := widget.NewLabel("")
	// 开启按钮
	check := widget.NewCheck("", func(b bool) {})
	check.Resize(fyne.NewSize(37, 37))
	check.Move(fyne.NewPos(313-47-39, 0))
	// 编辑按钮
	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
		createSshContent(w, groupName.Text).Show()
	})
	edit.Resize(fyne.NewSize(37, 37))
	edit.Move(fyne.NewPos(313-47, 0))
	return container.NewWithoutLayout(groupName, check, edit)
}

// 更新分组的元素
func updateGroupItem(item binding.DataItem, object fyne.CanvasObject) {
	con := object.(*fyne.Container)
	label := con.Objects[0].(*widget.Label)
	s, _ := item.(binding.String).Get()
	// [ssh名称][是否选中][分组坐标]
	ss := strings.Split(s, ",")
	label.SetText(ss[0])
	check := con.Objects[1].(*widget.Check)
	onOff := ss[1] == "true"
	check.SetChecked(onOff)
	i, _ := strconv.Atoi(ss[2])
	// 改变开启状态
	check.OnChanged = func(sel bool) {
		(*config.Yml.Log.Group)[i].OnOff = sel
	}
}

// 运行或暂停
func runOrStop(button *widget.Button, a *fyne.App, w fyne.Window) {
	runStatus = !runStatus
	if runStatus {
		config.SaveLocalCache()
		service.CreateSshClient()
		buildLogWindow(a)
		if len(logWin) > 0 {
			button.SetIcon(theme.MediaStopIcon())
		} else {
			dialogInfo("没有任何打开的链接", "需要添加或打开任意ssh链接", w)
		}
	} else {
		button.SetIcon(theme.MediaPlayIcon())
		service.CloseSshClient()
		closeLogWin()
	}
}

// 构建日志窗口
func buildLogWindow(a *fyne.App) {
	// 分组是否已经创建的map判断
	groupOnMap := make(map[string]bool, service.GroupNum)
	// 构建窗口对象
	for _, group := range *config.Yml.Log.Group {
		if !group.OnOff {
			continue
		}
		for _, ssh := range *group.Ssh {
			// 判断 服务器是否开启 是否已经创建
			if !ssh.OnOff || groupOnMap[group.Name] {
				continue
			}
			groupOnMap[group.Name] = true
			w := (*a).NewWindow(group.Name)
			w.SetContent(logView(group.Name))
			w.Resize(fyne.NewSize(1024, 648))
			w.CenterOnScreen()
			w.SetFixedSize(true)
			// 显示窗口
			w.Show()
			logWin = append(logWin, w)
		}
	}
}

// 关闭日志窗口
func closeLogWin() {
	for _, window := range logWin {
		window.Close()
	}
}

// 主题切换
func themeChange(btn *widget.Button) {
	if CurrentTheme {
		fyne.CurrentApp().Settings().SetTheme(theme.DarkTheme())
		btn.SetIcon(icon.CircleWhite)
	} else {
		fyne.CurrentApp().Settings().SetTheme(theme.LightTheme())
		btn.SetIcon(icon.CircleBlack)
	}
}

// 获取适合当前主题的图标
func getThemeIcon() *fyne.StaticResource {
	if CurrentTheme {
		return icon.CircleBlack
	}
	return icon.CircleWhite
}

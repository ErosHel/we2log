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
)

// 创建ssh管理页面dialog
func createSshContent(w fyne.Window, group string) dialog.Dialog {
	// 添加按钮
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	addBtn.Resize(fyne.NewSize(37, 37))

	// 分割线
	separator := widget.NewSeparator()
	separator.Resize(fyne.NewSize(275, 1))
	separator.Move(fyne.NewPos(0, 39))

	// 列表
	list := getGroupSshBindList(group)
	var groupList *widget.List
	groupList = widget.NewListWithData(list, createSshItem, func(item binding.DataItem, object fyne.CanvasObject) {
		updateSshItem(item, object, w, groupList)
	})
	groupList.Move(fyne.NewPos(0, 40))
	groupList.Resize(fyne.NewSize(280, 442))
	addBtn.OnTapped = func() {
		conf := &config.Ssh{}
		createSshEditContent(w, group, conf, func(name string) {
			g := *config.Yml.Log.Group
			// 获取分组以及ssh链接坐标
			gIndex := 0
			sshIndex := 0
			for i, g2 := range g {
				if g2.Name == group {
					gIndex = i
					// 增加ssh链接
					*g2.Ssh = append(*g2.Ssh, *conf)
					sshIndex = len(*g2.Ssh) - 1
					break
				}
			}
			// [ssh名称][分组坐标][链接坐标]
			_ = list.Append(fmt.Sprintf("%s,%d,%d", name, gIndex, sshIndex))
			groupList.Refresh()
			config.SaveLocalCache()
		}).Show()
	}
	// 弹窗内容
	content := container.NewWithoutLayout(addBtn, separator, groupList)
	// 创建弹窗
	custom := dialog.NewCustom(fmt.Sprintf("SSH链接管理 - %s", group), "返回", content, w)
	custom.Resize(fyne.NewSize(313, 480))
	return custom
}

// 获取分组下的ssh链接绑定列表
func getGroupSshBindList(group string) binding.StringList {
	groups := *config.Yml.Log.Group
	for i, g := range groups {
		if g.Name == group {
			ss := make([]string, len(*g.Ssh))
			for si, ssh := range *g.Ssh {
				// [ssh名称][分组坐标][链接坐标]
				ss[si] = fmt.Sprintf("%s,%d,%d", ssh.Name, i, si)
			}
			list := binding.NewStringList()
			_ = list.Set(ss)
			return list
		}
	}
	return nil
}

// 创建ssh item
func createSshItem() fyne.CanvasObject {
	// 分组名称
	sshName := widget.NewLabel("")
	// 开启按钮
	check := widget.NewCheck("", func(b bool) {})
	check.Resize(fyne.NewSize(37, 37))
	check.Move(fyne.NewPos(280-60-39, 0))
	// 编辑按钮
	edit := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {})
	edit.Resize(fyne.NewSize(37, 37))
	edit.Move(fyne.NewPos(280-60, 0))
	return container.NewWithoutLayout(sshName, check, edit)
}

// 更新ssh item
func updateSshItem(item binding.DataItem, object fyne.CanvasObject, w fyne.Window, list *widget.List) {
	con := object.(*fyne.Container)
	label := con.Objects[0].(*widget.Label)
	bStr := item.(binding.String)
	s, _ := bStr.Get()
	// [ssh名称][分组坐标][链接坐标]
	ss := strings.Split(s, ",")
	label.SetText(ss[0])
	check := con.Objects[1].(*widget.Check)
	i, _ := strconv.Atoi(ss[1])
	si, _ := strconv.Atoi(ss[2])
	onOff := (*(*config.Yml.Log.Group)[i].Ssh)[si].OnOff
	check.SetChecked(onOff)
	// 改变开启状态
	check.OnChanged = func(sel bool) {
		(*(*config.Yml.Log.Group)[i].Ssh)[si].OnOff = sel
	}
	btn := con.Objects[2].(*widget.Button)
	btn.OnTapped = createSshEditContent(w, ss[0], &(*(*config.Yml.Log.Group)[i].Ssh)[si], func(name string) {
		// [ssh名称][分组坐标][链接坐标]
		_ = bStr.Set(fmt.Sprintf("%s,%d,%d", name, i, si))
		list.Refresh()
	}).Show
}

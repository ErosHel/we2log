package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strings"
	"we2log/config"
)

// 创建ssh编辑页面dialog
func createSshEditContent(w fyne.Window, group string, conf *config.Ssh, confirmCallBack func(string)) dialog.Dialog {
	// 名称栏
	nameEntry := widget.NewEntry()
	nameEntry.SetText(conf.Name)
	nameEntry.SetPlaceHolder("服务器名称，用做识别")
	nameItem := widget.NewFormItem("名称", nameEntry)
	// 地址栏
	siteEntry := widget.NewEntry()
	siteEntry.SetText(conf.Host)
	siteEntry.SetPlaceHolder("服务器地址")
	siteEntry.Resize(fyne.NewSize(130, 37))
	// 端口栏
	portEntry := widget.NewEntry()
	if conf.Port == "" {
		portEntry.SetText("22")
	} else {
		portEntry.SetText(conf.Port)
	}
	portEntry.SetPlaceHolder("端口")
	portEntry.Resize(fyne.NewSize(60, 37))
	portEntry.Move(fyne.NewPos(133, 0))
	// 地址+端口
	siteLayout := container.NewWithoutLayout(siteEntry, portEntry)
	siteItem := widget.NewFormItem("地址", siteLayout)
	// 用户名栏
	userEntry := widget.NewEntry()
	userEntry.SetText(conf.Username)
	userEntry.SetPlaceHolder("服务器登陆用户")
	userItem := widget.NewFormItem("用户", userEntry)
	// 密码
	pwEntry := widget.NewPasswordEntry()
	pwEntry.SetText(conf.Password)
	if conf.PwType == 1 {
		pwEntry.Disabled()
	}
	pwEntry.SetPlaceHolder("服务器登陆密码")
	pwItem := widget.NewFormItem("密码", pwEntry)
	// 私钥路径栏
	priPathEntry := widget.NewEntry()
	priPathEntry.SetText(conf.PriKeyPath)
	priPathEntry.SetPlaceHolder("私钥路径(可空)")
	priPathEntry.Resize(fyne.NewSize(150, 37))
	// 私钥打开路径
	priPathOpen := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if closer != nil {
			priPathEntry.SetText(strings.ReplaceAll(closer.URI().String(), "file://", ""))
			_ = closer.Close()
		}
	}, w)
	priPathOpen.Resize(fyne.NewSize(600, 500))
	// 私钥选择按钮
	priSelectBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), priPathOpen.Show)
	priSelectBtn.Resize(fyne.NewSize(37, 37))
	priSelectBtn.Move(fyne.NewPos(153, 0))
	priPathLayout := container.NewWithoutLayout(priPathEntry, priSelectBtn)
	priPathItem := widget.NewFormItem("私钥路径", priPathLayout)
	if conf.PwType == 0 {
		priPathEntry.Disable()
		priSelectBtn.Disable()
	}
	// 密码类型
	pwType := widget.NewSelect([]string{"密码", "私钥"}, func(s string) {
		if s == "密码" {
			priPathEntry.Disable()
			priSelectBtn.Disable()
			pwEntry.Enable()
		} else {
			pwEntry.Disable()
			priPathEntry.Enable()
			priSelectBtn.Enable()
		}
	})
	pwType.SetSelectedIndex(conf.PwType)
	pwTypeItem := widget.NewFormItem("密码类型", pwType)
	// 日志路径
	logPathEntry := widget.NewEntry()
	logPathEntry.SetText(conf.LogPath)
	logPathEntry.SetPlaceHolder("服务器日志文件路径")
	logPathItem := widget.NewFormItem("日志路径", logPathEntry)
	formItem := []*widget.FormItem{nameItem, siteItem, userItem, pwItem, priPathItem, pwTypeItem, logPathItem}
	form := dialog.NewForm(fmt.Sprintf("SSH编辑页面 - %s", group), "保存", "取消", formItem, func(confirm bool) {
		if confirm {
			// 参数校验
			if nameEntry.Text == "" {
				dialogInfo("参数错误", "服务器名称不能为空", w)
				return
			}
			if siteEntry.Text == "" {
				dialogInfo("参数错误", "服务器地址不能为空", w)
				return
			}
			if portEntry.Text == "" {
				dialogInfo("参数错误", "服务器端口不能为空", w)
				return
			}
			if userEntry.Text == "" {
				dialogInfo("参数错误", "服务器用户不能为空", w)
				return
			}
			if pwType.SelectedIndex() == 0 && pwEntry.Text == "" {
				dialogInfo("参数错误", "密码类型时，密码不能为空", w)
				return
			}
			if logPathEntry.Text == "" {
				dialogInfo("参数错误", "日志路径不能为空", w)
				return
			}
			// 参数赋值
			conf.Name = nameEntry.Text
			conf.Host = siteEntry.Text
			conf.Port = portEntry.Text
			conf.Username = userEntry.Text
			conf.Password = pwEntry.Text
			conf.PriKeyPath = priPathEntry.Text
			conf.PwType = pwType.SelectedIndex()
			conf.LogPath = logPathEntry.Text
			confirmCallBack(nameEntry.Text)
		}
	}, w)
	form.Resize(fyne.NewSize(313, 450))
	return form
}

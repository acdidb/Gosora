// +build !no_templategen

// Code generated by Gosora. More below:
/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */
package main
import "net/http"
import "./common"
import "strconv"

// nolint
func init() {
	common.Template_profile_handle = Template_profile
	common.Ctemplates = append(common.Ctemplates,"profile")
	common.TmplPtrMap["profile"] = &common.Template_profile_handle
	common.TmplPtrMap["o_profile"] = Template_profile
}

// nolint
func Template_profile(tmpl_profile_vars common.ProfilePage, w http.ResponseWriter) error {
w.Write(header_0)
w.Write([]byte(tmpl_profile_vars.Title))
w.Write(header_1)
w.Write([]byte(tmpl_profile_vars.Header.Site.Name))
w.Write(header_2)
w.Write([]byte(tmpl_profile_vars.Header.Theme.Name))
w.Write(header_3)
if len(tmpl_profile_vars.Header.Stylesheets) != 0 {
for _, item := range tmpl_profile_vars.Header.Stylesheets {
w.Write(header_4)
w.Write([]byte(item))
w.Write(header_5)
}
}
w.Write(header_6)
if len(tmpl_profile_vars.Header.Scripts) != 0 {
for _, item := range tmpl_profile_vars.Header.Scripts {
w.Write(header_7)
w.Write([]byte(item))
w.Write(header_8)
}
}
w.Write(header_9)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(header_10)
w.Write([]byte(tmpl_profile_vars.Header.Site.URL))
w.Write(header_11)
if tmpl_profile_vars.Header.MetaDesc != "" {
w.Write(header_12)
w.Write([]byte(tmpl_profile_vars.Header.MetaDesc))
w.Write(header_13)
}
w.Write(header_14)
if !tmpl_profile_vars.CurrentUser.IsSuperMod {
w.Write(header_15)
}
w.Write(header_16)
w.Write(menu_0)
w.Write(menu_1)
w.Write([]byte(tmpl_profile_vars.Header.Site.ShortName))
w.Write(menu_2)
if tmpl_profile_vars.CurrentUser.Loggedin {
w.Write(menu_3)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Link))
w.Write(menu_4)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(menu_5)
} else {
w.Write(menu_6)
}
w.Write(menu_7)
w.Write(header_17)
if tmpl_profile_vars.Header.Widgets.RightSidebar != "" {
w.Write(header_18)
}
w.Write(header_19)
if len(tmpl_profile_vars.Header.NoticeList) != 0 {
for _, item := range tmpl_profile_vars.Header.NoticeList {
w.Write(header_20)
w.Write([]byte(item))
w.Write(header_21)
}
}
w.Write(profile_0)
w.Write([]byte(tmpl_profile_vars.ProfileOwner.Avatar))
w.Write(profile_1)
w.Write([]byte(tmpl_profile_vars.ProfileOwner.Name))
w.Write(profile_2)
w.Write([]byte(tmpl_profile_vars.ProfileOwner.Name))
w.Write(profile_3)
w.Write([]byte(tmpl_profile_vars.ProfileOwner.Name))
w.Write(profile_4)
if tmpl_profile_vars.ProfileOwner.Tag != "" {
w.Write(profile_5)
w.Write([]byte(tmpl_profile_vars.ProfileOwner.Tag))
w.Write(profile_6)
}
w.Write(profile_7)
if !tmpl_profile_vars.CurrentUser.Loggedin {
w.Write(profile_8)
} else {
w.Write(profile_9)
if tmpl_profile_vars.CurrentUser.IsSuperMod && !tmpl_profile_vars.ProfileOwner.IsSuperMod {
w.Write(profile_10)
if tmpl_profile_vars.ProfileOwner.IsBanned {
w.Write(profile_11)
w.Write([]byte(strconv.Itoa(tmpl_profile_vars.ProfileOwner.ID)))
w.Write(profile_12)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_13)
} else {
w.Write(profile_14)
}
w.Write(profile_15)
}
w.Write(profile_16)
w.Write([]byte(strconv.Itoa(tmpl_profile_vars.ProfileOwner.ID)))
w.Write(profile_17)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_18)
}
w.Write(profile_19)
if tmpl_profile_vars.CurrentUser.Perms.BanUsers {
w.Write(profile_20)
w.Write([]byte(strconv.Itoa(tmpl_profile_vars.ProfileOwner.ID)))
w.Write(profile_21)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_22)
w.Write(profile_23)
}
w.Write(profile_24)
if tmpl_profile_vars.Header.Theme.BgAvatars {
if len(tmpl_profile_vars.ItemList) != 0 {
for _, item := range tmpl_profile_vars.ItemList {
w.Write(profile_comments_row_0)
w.Write([]byte(item.ClassName))
w.Write(profile_comments_row_1)
w.Write([]byte(item.Avatar))
w.Write(profile_comments_row_2)
if item.ContentLines <= 5 {
w.Write(profile_comments_row_3)
}
w.Write(profile_comments_row_4)
w.Write([]byte(item.ContentHtml))
w.Write(profile_comments_row_5)
w.Write([]byte(item.UserLink))
w.Write(profile_comments_row_6)
w.Write([]byte(item.CreatedByName))
w.Write(profile_comments_row_7)
if tmpl_profile_vars.CurrentUser.IsMod {
w.Write(profile_comments_row_8)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(profile_comments_row_9)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_comments_row_10)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(profile_comments_row_11)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_comments_row_12)
}
w.Write(profile_comments_row_13)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(profile_comments_row_14)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_comments_row_15)
if item.Tag != "" {
w.Write(profile_comments_row_16)
w.Write([]byte(item.Tag))
w.Write(profile_comments_row_17)
}
w.Write(profile_comments_row_18)
}
}
} else {
if len(tmpl_profile_vars.ItemList) != 0 {
for _, item := range tmpl_profile_vars.ItemList {
w.Write(profile_comments_row_19)
w.Write([]byte(item.ClassName))
w.Write(profile_comments_row_20)
w.Write([]byte(item.Avatar))
w.Write(profile_comments_row_21)
w.Write([]byte(item.CreatedByName))
w.Write(profile_comments_row_22)
w.Write([]byte(item.CreatedByName))
w.Write(profile_comments_row_23)
w.Write([]byte(item.UserLink))
w.Write(profile_comments_row_24)
w.Write([]byte(item.CreatedByName))
w.Write(profile_comments_row_25)
if item.Tag != "" {
w.Write(profile_comments_row_26)
w.Write([]byte(item.Tag))
w.Write(profile_comments_row_27)
}
w.Write(profile_comments_row_28)
if tmpl_profile_vars.CurrentUser.IsMod {
w.Write(profile_comments_row_29)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(profile_comments_row_30)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_comments_row_31)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(profile_comments_row_32)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_comments_row_33)
}
w.Write(profile_comments_row_34)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(profile_comments_row_35)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_comments_row_36)
w.Write([]byte(item.ContentHtml))
w.Write(profile_comments_row_37)
}
}
}
w.Write(profile_25)
if !tmpl_profile_vars.CurrentUser.IsBanned {
w.Write(profile_26)
w.Write([]byte(tmpl_profile_vars.CurrentUser.Session))
w.Write(profile_27)
w.Write([]byte(strconv.Itoa(tmpl_profile_vars.ProfileOwner.ID)))
w.Write(profile_28)
}
w.Write(profile_29)
w.Write(profile_30)
w.Write(footer_0)
w.Write([]byte(common.BuildWidget("footer",tmpl_profile_vars.Header)))
w.Write(footer_1)
if len(tmpl_profile_vars.Header.Themes) != 0 {
for _, item := range tmpl_profile_vars.Header.Themes {
if !item.HideFromThemes {
w.Write(footer_2)
w.Write([]byte(item.Name))
w.Write(footer_3)
if tmpl_profile_vars.Header.Theme.Name == item.Name {
w.Write(footer_4)
}
w.Write(footer_5)
w.Write([]byte(item.FriendlyName))
w.Write(footer_6)
}
}
}
w.Write(footer_7)
w.Write([]byte(common.BuildWidget("rightSidebar",tmpl_profile_vars.Header)))
w.Write(footer_8)
	return nil
}

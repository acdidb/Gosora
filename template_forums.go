/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */
package main
import "strconv"
import "io"

func init() {
	template_forums_handle = template_forums
	//o_template_forums_handle = template_forums
	ctemplates = append(ctemplates,"forums")
	tmpl_ptr_map["forums"] = &template_forums_handle
	tmpl_ptr_map["o_forums"] = template_forums
}

func template_forums(tmpl_forums_vars ForumsPage, w io.Writer) {
w.Write(header_0)
w.Write([]byte(tmpl_forums_vars.Title))
w.Write(header_1)
w.Write([]byte(tmpl_forums_vars.CurrentUser.Session))
w.Write(header_2)
w.Write(menu_0)
if tmpl_forums_vars.CurrentUser.Loggedin {
w.Write(menu_1)
w.Write([]byte(strconv.Itoa(tmpl_forums_vars.CurrentUser.ID)))
w.Write(menu_2)
if tmpl_forums_vars.CurrentUser.Is_Super_Mod {
w.Write(menu_3)
}
w.Write(menu_4)
w.Write([]byte(tmpl_forums_vars.CurrentUser.Session))
w.Write(menu_5)
} else {
w.Write(menu_6)
}
w.Write(menu_7)
w.Write(header_3)
if len(tmpl_forums_vars.NoticeList) != 0 {
for _, item := range tmpl_forums_vars.NoticeList {
w.Write(header_4)
w.Write([]byte(item))
w.Write(header_5)
}
}
w.Write(forums_0)
if len(tmpl_forums_vars.ItemList) != 0 {
for _, item := range tmpl_forums_vars.ItemList {
w.Write(forums_1)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(forums_2)
w.Write([]byte(item.Name))
w.Write(forums_3)
w.Write([]byte(strconv.Itoa(item.LastTopicID)))
w.Write(forums_4)
w.Write([]byte(item.LastTopic))
w.Write(forums_5)
w.Write([]byte(item.LastTopicTime))
w.Write(forums_6)
}
} else {
w.Write(forums_7)
}
w.Write(forums_8)
w.Write(footer_0)
}
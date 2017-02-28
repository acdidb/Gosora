/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */
package main
import "io"
import "strconv"

func init() {
	template_topics_handle = template_topics
	//o_template_topics_handle = template_topics
	ctemplates = append(ctemplates,"topics")
	tmpl_ptr_map["topics"] = &template_topics_handle
	tmpl_ptr_map["o_topics"] = template_topics
}

func template_topics(tmpl_topics_vars TopicsPage, w io.Writer) {
w.Write(header_0)
w.Write([]byte(tmpl_topics_vars.Title))
w.Write(header_1)
w.Write([]byte(tmpl_topics_vars.CurrentUser.Session))
w.Write(header_2)
w.Write(menu_0)
if tmpl_topics_vars.CurrentUser.Loggedin {
w.Write(menu_1)
w.Write([]byte(strconv.Itoa(tmpl_topics_vars.CurrentUser.ID)))
w.Write(menu_2)
if tmpl_topics_vars.CurrentUser.Is_Super_Mod {
w.Write(menu_3)
}
w.Write(menu_4)
w.Write([]byte(tmpl_topics_vars.CurrentUser.Session))
w.Write(menu_5)
} else {
w.Write(menu_6)
}
w.Write(menu_7)
if !tmpl_topics_vars.CurrentUser.Loggedin {
w.Write(menu_8)
}
w.Write(menu_9)
w.Write(header_3)
if len(tmpl_topics_vars.NoticeList) != 0 {
for _, item := range tmpl_topics_vars.NoticeList {
w.Write(header_4)
w.Write([]byte(item))
w.Write(header_5)
}
}
w.Write(topics_0)
if len(tmpl_topics_vars.ItemList) != 0 {
for _, item := range tmpl_topics_vars.ItemList {
w.Write(topics_1)
if item.Avatar != "" {
w.Write(topics_2)
w.Write([]byte(item.Avatar))
w.Write(topics_3)
}
if item.Sticky {
w.Write(topics_4)
} else {
if item.Is_Closed {
w.Write(topics_5)
}
}
w.Write(topics_6)
w.Write([]byte(strconv.Itoa(item.ID)))
w.Write(topics_7)
w.Write([]byte(item.Title))
w.Write(topics_8)
if item.ForumName != "" {
w.Write(topics_9)
w.Write([]byte(strconv.Itoa(item.ParentID)))
w.Write(topics_10)
w.Write([]byte(item.ForumName))
w.Write(topics_11)
}
if item.Is_Closed {
w.Write(topics_12)
}
w.Write(topics_13)
w.Write([]byte(item.LastReplyAt))
w.Write(topics_14)
}
} else {
w.Write(topics_15)
if tmpl_topics_vars.CurrentUser.Perms.CreateTopic {
w.Write(topics_16)
}
w.Write(topics_17)
}
w.Write(topics_18)
w.Write(footer_0)
}

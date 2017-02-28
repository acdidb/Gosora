package main
//import "log"
//import "fmt"
import "bytes"
//import "strings"
import "strconv"
import "regexp"
import "time"
import "math/rand"

var random *rand.Rand
var bbcode_invalid_url []byte
var bbcode_invalid_number []byte
var bbcode_missing_tag []byte
var bbcode_url_open []byte
var bbcode_url_open2 []byte
var bbcode_url_close []byte
var bbcode_space_gap []byte

var bbcode_bold *regexp.Regexp
var bbcode_italic *regexp.Regexp
var bbcode_underline *regexp.Regexp
var bbcode_strikethrough *regexp.Regexp
var bbcode_url *regexp.Regexp
var bbcode_url_label *regexp.Regexp

func init() {
	plugins["bbcode"] = NewPlugin("bbcode","BBCode","Azareal","http://github.com/Azareal","","","",init_bbcode,nil,deactivate_bbcode)
}

func init_bbcode() {
	//plugins["bbcode"].AddHook("parse_assign", bbcode_parse_without_code)
	plugins["bbcode"].AddHook("parse_assign", bbcode_full_parse)
	
	bbcode_invalid_url = []byte("<span style='color: red;'>[Invalid URL]</span>")
	bbcode_invalid_number = []byte("<span style='color: red;'>[Invalid Number]</span>")
	bbcode_missing_tag = []byte("<span style='color: red;'>[Missing Tag]</span>")
	bbcode_url_open = []byte("<a href='")
	bbcode_url_open2 = []byte("'>")
	bbcode_url_close = []byte("</a>")
	bbcode_space_gap = []byte("          ")
	
	bbcode_bold = regexp.MustCompile(`(?s)\[b\](.*)\[/b\]`)
	bbcode_italic = regexp.MustCompile(`(?s)\[i\](.*)\[/i\]`)
	bbcode_underline = regexp.MustCompile(`(?s)\[u\](.*)\[/u\]`)
	bbcode_strikethrough = regexp.MustCompile(`(?s)\[s\](.*)\[/s\]`)
	urlpattern := `(http|https|ftp|mailto*)(:??)\/\/([\.a-zA-Z\/]+)`
	bbcode_url = regexp.MustCompile(`\[url\]` + urlpattern + `\[/url\]`)
	bbcode_url_label = regexp.MustCompile(`(?s)\[url=` + urlpattern + `\](.*)\[/url\]`)
	
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func deactivate_bbcode() {
	//plugins["bbcode"].RemoveHook("parse_assign", bbcode_parse_without_code)
	plugins["bbcode"].RemoveHook("parse_assign", bbcode_full_parse)
}

func bbcode_regex_parse(data interface{}) interface{} {
	msg := data.(string)
	msg = bbcode_bold.ReplaceAllString(msg,"<b>$1</b>")
	msg = bbcode_italic.ReplaceAllString(msg,"<i>$1</i>")
	msg = bbcode_underline.ReplaceAllString(msg,"<u>$1</u>")
	msg = bbcode_strikethrough.ReplaceAllString(msg,"<s>$1</s>")
	msg = bbcode_url.ReplaceAllString(msg,"<a href=\"$1$2//$3\" rel=\"nofollow\">$1$2//$3</i>")
	msg = bbcode_url_label.ReplaceAllString(msg,"<a href=\"$1$2//$3\" rel=\"nofollow\">$4</i>")
	return msg
}

// Only does the simple BBCode like [u], [b], [i] and [s]
func bbcode_simple_parse(data interface{}) interface{} {
	msg := data.(string)
	msgbytes := []byte(msg)
	has_u := false
	has_b := false
	has_i := false
	has_s := false
	for i := 0; (i + 2) < len(msgbytes); i++ {
		if msgbytes[i] == '[' && msgbytes[i + 2] == ']' {
			if msgbytes[i + 1] == 'b' {
				msgbytes[i] = '<'
				msgbytes[i + 2] = '>'
				has_b = true
			} else if msgbytes[i + 1] == 'i' {
				msgbytes[i] = '<'
				msgbytes[i + 2] = '>'
				has_i = true
			} else if msgbytes[i + 1] == 'u' {
				msgbytes[i] = '<'
				msgbytes[i + 2] = '>'
				has_u = true
			} else if msgbytes[i + 1] == 's' {
				msgbytes[i] = '<'
				msgbytes[i + 2] = '>'
				has_s = true
			}
			i += 2
		}
	}
	
	// There's an unclosed tag in there x.x
	if has_i || has_u || has_b || has_s {
		closer := []byte("</u></i></b></s>")
		msgbytes = append(msgbytes, closer...)
	}
	return string(msgbytes)
}

// Here for benchmarking purposes. Might add a plugin setting for disabling [code] as it has it's paws everywhere
func bbcode_parse_without_code(data interface{}) interface{} {
	msg := data.(string)
	msgbytes := []byte(msg)
	has_u := false
	has_b := false
	has_i := false
	has_s := false
	complex_bbc := false
	for i := 0; (i + 3) < len(msgbytes); i++ {
		if msgbytes[i] == '[' {
			if msgbytes[i + 2] != ']' {
				if msgbytes[i + 1] == '/' {
					if msgbytes[i + 3] == ']' {
						if msgbytes[i + 2] == 'b' {
							msgbytes[i] = '<'
							msgbytes[i + 3] = '>'
							has_b = false
						} else if msgbytes[i + 2] == 'i' {
							msgbytes[i] = '<'
							msgbytes[i + 3] = '>'
							has_i = false
						} else if msgbytes[i + 2] == 'u' {
							msgbytes[i] = '<'
							msgbytes[i + 3] = '>'
							has_u = false
						} else if msgbytes[i + 2] == 's' {
							msgbytes[i] = '<'
							msgbytes[i + 3] = '>'
							has_s = false
						}
						i += 3
					} else {
						complex_bbc = true
					}
				} else {
					complex_bbc = true
				}
			} else {
				if msgbytes[i + 1] == 'b' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_b = true
				} else if msgbytes[i + 1] == 'i' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_i = true
				} else if msgbytes[i + 1] == 'u' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_u = true
				} else if msgbytes[i + 1] == 's' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_s = true
				}
				i += 2
			}
		}
	}
	
	// There's an unclosed tag in there x.x
	if has_i || has_u || has_b || has_s {
		closer := []byte("</u></i></b></s>")
		msgbytes = append(msgbytes, closer...)
	}
	
	// Copy the new complex parser over once the rough edges have been smoothed over
	if complex_bbc {
		msg = bbcode_url.ReplaceAllString(msg,"<a href=\"$1$2//$3\" rel=\"nofollow\">$1$2//$3</i>")
		msg = bbcode_url_label.ReplaceAllString(msg,"<a href=\"$1$2//$3\" rel=\"nofollow\">$4</i>")
	}
	return string(msgbytes)
}

// Does every type of BBCode
func bbcode_full_parse(data interface{}) interface{} {
	msg := data.(string)
	msgbytes := []byte(msg)
	has_u := false
	has_b := false
	has_i := false
	has_s := false
	has_c := false
	complex_bbc := false
	msgbytes = append(msgbytes,bbcode_space_gap...)
	for i := 0; i < len(msgbytes); i++ {
		if msgbytes[i] == '[' {
			if msgbytes[i + 2] != ']' {
				if msgbytes[i + 1] == '/' {
					if msgbytes[i + 3] == ']' {
						if !has_c {
							if msgbytes[i + 2] == 'b' {
								msgbytes[i] = '<'
								msgbytes[i + 3] = '>'
								has_b = false
							} else if msgbytes[i + 2] == 'i' {
								msgbytes[i] = '<'
								msgbytes[i + 3] = '>'
								has_i = false
							} else if msgbytes[i + 2] == 'u' {
								msgbytes[i] = '<'
								msgbytes[i + 3] = '>'
								has_u = false
							} else if msgbytes[i + 2] == 's' {
								msgbytes[i] = '<'
								msgbytes[i + 3] = '>'
								has_s = false
							}
							i += 3
						}
					} else {
						if msgbytes[i+2] == 'c' && msgbytes[i+3] == 'o' && msgbytes[i+4] == 'd' && msgbytes[i+5] == 'e' && msgbytes[i+6] == ']' {
							has_c = false
							i += 7
						}
						//if msglen >= (i+6) {
						//	fmt.Println("boo")
						//	fmt.Println(msglen)
						//	fmt.Println(i+6)
						//	fmt.Println(string(msgbytes[i:i+6]))
						//}
						complex_bbc = true
					}
				} else {
					if msgbytes[i+1] == 'c' && msgbytes[i+2] == 'o' && msgbytes[i+3] == 'd' && msgbytes[i+4] == 'e' && msgbytes[i+5] == ']' {
						has_c = true
						i += 6
					}
					//if msglen >= (i+5) {
					//	fmt.Println("boo2")
					//	fmt.Println(string(msgbytes[i:i+5]))
					//}
					complex_bbc = true
				}
			} else if !has_c {
				if msgbytes[i + 1] == 'b' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_b = true
				} else if msgbytes[i + 1] == 'i' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_i = true
				} else if msgbytes[i + 1] == 'u' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_u = true
				} else if msgbytes[i + 1] == 's' {
					msgbytes[i] = '<'
					msgbytes[i + 2] = '>'
					has_s = true
				}
				i += 2
			}
		}
	}
	
	// There's an unclosed tag in there x.x
	if has_i || has_u || has_b || has_s {
		closer := []byte("</u></i></b></s>")
		msgbytes = append(msgbytes, closer...)
	}
	
	if complex_bbc {
		i := 0
		var start int
		var lastTag int
		outbytes := make([]byte, len(msgbytes))
		//fmt.Println(string(msgbytes))
		for ; i < len(msgbytes); i++ {
			MainLoop:
			if msgbytes[i] == '[' {
				OuterComplex:
				if msgbytes[i + 1] == 'u' {
					if msgbytes[i+2] == 'r' && msgbytes[i+3] == 'l' && msgbytes[i+4] == ']' {
						outbytes = append(outbytes, msgbytes[lastTag:i]...)
						start = i + 5
						i = start
						if msgbytes[i] == 'h' {
							if msgbytes[i+1] == 't' && msgbytes[i+2] == 't' && msgbytes[i+3] == 'p' {
								if bytes.Equal(msgbytes[i + 4:i + 7],[]byte("s://")) {
									i += 7
								} else if msgbytes[i+4] == ':' && msgbytes[i+5] == '/' && msgbytes[i+6] == '/' {
									i += 6
								} else {
									outbytes = append(outbytes, bbcode_invalid_url...)
									continue
								}
							}
						} else if msgbytes[i] == 'f' {
							if bytes.Equal(msgbytes[i+1:i+5],[]byte("tp://")) {
								i += 5
							}
						}
						
						for ;; i++ {
							if len(msgbytes) < (i + 10) {
								//fmt.Println(msglen)
								//fmt.Println(i+6)
								outbytes = append(outbytes, bbcode_missing_tag...)
								goto OuterComplex
							} else if msgbytes[i] == '[' {
								if !bytes.Equal(msgbytes[i+1:i+6],[]byte("/url]")) {
									//log.Print("Not the URL closing tag!")
									//fmt.Println(msgbytes[i + 1:i + 6])
									outbytes = append(outbytes, bbcode_missing_tag...)
									goto OuterComplex
								}
								break
							} else if msgbytes[i] != '\\' && msgbytes[i] != '_' && !(msgbytes[i] > 44 && msgbytes[i] < 58) && !(msgbytes[i] > 64 && msgbytes[i] < 91) && !(msgbytes[i] > 96 && msgbytes[i] < 123) {
								outbytes = append(outbytes, bbcode_invalid_url...)
								//log.Print("Weird character")
								//fmt.Println(msgbytes[i])
								goto MainLoop
							}
						}
						outbytes = append(outbytes, bbcode_url_open...)
						outbytes = append(outbytes, msgbytes[start:i]...)
						outbytes = append(outbytes, bbcode_url_open2...)
						outbytes = append(outbytes, msgbytes[start:i]...)
						outbytes = append(outbytes, bbcode_url_close...)
						i += 6
						lastTag = i
					}
				} else if msgbytes[i + 1] == 'r' {
					if bytes.Equal(msgbytes[i+2:i+6],[]byte("and]")) {
						outbytes = append(outbytes, msgbytes[lastTag:i]...)
						start = i + 6
						i = start
						for ;; i++ {
							if msgbytes[i] == '[' {
								if !bytes.Equal(msgbytes[i+1:i+7],[]byte("/rand]")) {
									outbytes = append(outbytes, bbcode_missing_tag...)
									goto OuterComplex
								}
								break
							} else if (len(msgbytes) - 1) < (i + 10) {
								outbytes = append(outbytes, bbcode_missing_tag...)
									goto OuterComplex
							}
						}
						
						number, err := strconv.ParseInt(string(msgbytes[start:i]),10,64)
						if err != nil {
							outbytes = append(outbytes, bbcode_invalid_number...)
							goto MainLoop
						}
						
						dat := []byte(strconv.FormatInt((random.Int63n(number)),10))
						outbytes = append(outbytes, dat...)
						//log.Print("Outputted the random number")
						i += 7
						lastTag = i
					}
				}
			}
		}
		//fmt.Println(string(outbytes))
		if lastTag != i {
			outbytes = append(outbytes, msgbytes[lastTag:]...)
		}
		if len(outbytes) != 0 {
			return string(outbytes[0:len(msgbytes) - 10])
		}
		
		msg = string(msgbytes[0:len(msgbytes) - 10])
		//msg = bbcode_url.ReplaceAllString(msg,"<a href=\"$1$2//$3\" rel=\"nofollow\">$1$2//$3</i>")
		msg = bbcode_url_label.ReplaceAllString(msg,"<a href=\"$1$2//$3\" rel=\"nofollow\">$4</i>")
		// Convert [code] into class="codequotes"
	} else {
		msg = string(msgbytes[0:len(msgbytes) - 10])
	}
	return msg
}

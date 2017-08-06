/* WIP Under Construction */
package main

import "log"
import "fmt"
//import "strings"
import "os"

var route_list []Route
var route_groups []RouteGroup

func main() {
	fmt.Println("Generating the router...")
	
	// Load all the routes...
	routes()
	
	var out string
	var fdata string = "// Code generated by. DO NOT EDIT.\n/* This file was automatically generated by the software. Please don't edit it as your changes may be overwritten at any moment. */\n"
	
	for _, route := range route_list {
		var end int
		if route.Path[len(route.Path)-1] == '/' {
			end = len(route.Path) - 1
		} else {
			end = len(route.Path) - 1
		}
		out += "\n\t\tcase \"" + route.Path[0:end] + "\":"
		if route.Before != "" {
			out += "\n\t\t\t" + route.Before
		}
		out += "\n\t\t\t" + route.Name + "(w,req,user"
		for _, item := range route.Vars {
			out += "," + item
		}
		out += ")\n\t\t\treturn"
	}
	
	for _, group := range route_groups {
		var end int
		if group.Path[len(group.Path)-1] == '/' {
			end = len(group.Path) - 1
		} else {
			end = len(group.Path) - 1
		}
		out += `
		case "` + group.Path[0:end] + `":
			switch(req.URL.Path) {`
		var default_route Route
		for _, route := range group.Routes {
			if group.Path == route.Path {
				default_route = route
				continue
			}
			
			out += "\n\t\t\t\tcase \"" + route.Path + "\":"
			if route.Before != "" {
				out += "\n\t\t\t\t\t" + route.Before
			}
			out += "\n\t\t\t\t\t" + route.Name + "(w,req,user"
			for _, item := range route.Vars {
				out += "," + item
			}
			out += ")\n\t\t\t\t\treturn"
		}
		
		if default_route.Name != "" {
			out += "\n\t\t\t\tdefault:"
			if default_route.Before != "" {
				out += "\n\t\t\t\t\t" + default_route.Before
			}
			out += "\n\t\t\t\t\t" + default_route.Name + "(w,req,user"
			for _, item := range default_route.Vars {
				out += ", " + item
			}
			out += ")\n\t\t\t\t\treturn"
		}
		out += "\n\t\t\t}"
	}
	
	fdata += `package main

import "fmt"
import "strings"
import "sync"
import "errors"
import "net/http"

var ErrNoRoute = errors.New("That route doesn't exist.")

type GenRouter struct {
	UploadHandler func(http.ResponseWriter, *http.Request)
	extra_routes map[string]func(http.ResponseWriter, *http.Request, User)
	
	sync.RWMutex
}

func NewGenRouter(uploads http.Handler) *GenRouter {
	return &GenRouter{
		UploadHandler: http.StripPrefix("/uploads/",uploads).ServeHTTP,
		extra_routes: make(map[string]func(http.ResponseWriter, *http.Request, User)),
	}
}

func (router *GenRouter) Handle(_ string, _ http.Handler) {
}

func (router *GenRouter) HandleFunc(pattern string, handle func(http.ResponseWriter, *http.Request, User)) {
	router.Lock()
	router.extra_routes[pattern] = handle
	router.Unlock()
}

func (router *GenRouter) RemoveFunc(pattern string) error {
	router.Lock()
	_, ok := router.extra_routes[pattern]
	if !ok {
		router.Unlock()
		return ErrNoRoute
	}
	delete(router.extra_routes,pattern)
	router.Unlock()
	return nil
}

func (router *GenRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//if req.URL.Path == "/" {
	//	default_route(w,req)
	//	return
	//}
	if req.URL.Path[0] != '/' {
		w.WriteHeader(405)
		w.Write([]byte(""))
		return
	}
	
	var prefix, extra_data string
	prefix = req.URL.Path[0:strings.IndexByte(req.URL.Path[1:],'/') + 1]
	if req.URL.Path[len(req.URL.Path) - 1] != '/' {
		extra_data = req.URL.Path[strings.LastIndexByte(req.URL.Path,'/') + 1:]
		req.URL.Path = req.URL.Path[:strings.LastIndexByte(req.URL.Path,'/') + 1]
	}
	
	if dev.SuperDebug {
		fmt.Println("before route_static")
		fmt.Println("prefix:",prefix)
		fmt.Println("req.URL.Path:",req.URL.Path)
		fmt.Println("extra_data:",extra_data)
	}
	
	if prefix == "/static" {
		req.URL.Path += extra_data
		route_static(w,req)
		return
	}
	
	if dev.SuperDebug {
		fmt.Println("before PreRoute")
	}
	
	// Deal with the session stuff, etc.
	user, ok := PreRoute(w,req)
	if !ok {
		return
	}
	
	if dev.SuperDebug {
		fmt.Println("after PreRoute")
	}
	
	switch(prefix) {` + out + `
		case "/uploads":
			if extra_data == "" {
				NotFound(w,req)
				return
			}
			req.URL.Path += extra_data
			router.UploadHandler(w,req)
			return
		case "":
			// Stop the favicons, robots.txt file, etc. resolving to the topics list
			// TO-DO: Add support for favicons and robots.txt files
			switch(extra_data) {
				case "robots.txt":
					route_robots_txt(w,req)
					return
			}
			
			if extra_data != "" {
				NotFound(w,req)
				return
			}
			config.DefaultRoute(w,req,user)
			return
		//default: NotFound(w,req)
	}
	
	// A fallback for the routes which haven't been converted to the new router yet or plugins
	router.RLock()
	handle, ok := router.extra_routes[req.URL.Path]
	router.RUnlock()
	
	if ok {
		req.URL.Path += extra_data
		handle(w,req,user)
		return
	}
	NotFound(w,req)
}
`
	write_file("./gen_router.go",fdata)
	fmt.Println("Successfully generated the router")
}

func write_file(name string, content string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}
	f.Sync()
	f.Close()
}

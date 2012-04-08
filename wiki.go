package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"regexp"
	"errors"
	"github.com/russross/blackfriday"
	"bytes"
	"strings"
	"os"
)

type Page struct {
    Title string
    Body  []byte
}

var wikiname = "m33w#Paste"
var css_ = []byte("#")
var default_ = []byte("#")

const lenPath = len("/view/")
//var templates = template.Must(template.ParseFiles("html/edit.html", "html/view.html"))
var titleValidator = regexp.MustCompile("^[a-zA-Z0-9_]+$")

func (p *Page) save() error {
    filename := "pages/" + p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := "pages/" + title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func loadLib(title string) ([]byte) {
    body, err := ioutil.ReadFile(title)
    if err != nil {
        return []byte("")
    }
    return body
}

/*****************
**     VIEWS    **
*****************/

func renderView(w http.ResponseWriter ,p *Page, r *http.Request, markdown bool) {
		// Get Config from Body
		var boder = bytes.Split(p.Body,[]byte("~![META]!~"))
		// Replace Body with showable body
		p.Body = boder[1]
		if len(boder) > 1 {
			// Now lets look at the config meta
			var confer = bytes.Split(boder[0],[]byte("::"))
			// Init some trash vars
			var public_view = false
			var public_edit = false
			// And korrigate them
			if bytes.Equal(confer[0],[]byte("true")) {
				public_view = true
			}
			if bytes.Equal(confer[1],[]byte("true")) {
				public_edit = true
			}
			// Get Password
			var password = strings.TrimSpace(string(confer[2]))

			// Print Default Body [HEAD PART]
			fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \n \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n" +
					"<html xmlns=\"http://www.w3.org/1999/xhtml\">\n" +
					"<head>\n" +
					"	<title>%s</title>\n" +
					"	<style type=\"text/css\">%s</style> \n " +
					"</head>\n" +
					"<body onload=\"document.getElementById('_password').focus();\">\n	<div id='header'>\n	<div class=\"right\"><a href='/'>%s</a></div>\n	<h1><a href=\"/view/%s\">%s</a></h1>\n	</div>\n", p.Title+" ~ "+wikiname, css_, wikiname, p.Title, p.Title)

			if public_view || strings.TrimSpace(r.FormValue("password")) == password {
				if markdown {
					p.Body = blackfriday.MarkdownBasic(p.Body)
				}
				fmt.Fprintf(w, "	<div id='content'>\n%s\n	</div>\n", p.Body)
				fmt.Fprintf(w, "	<form action='/edit/%s' >",p.Title)
				if public_edit {
					fmt.Fprintf(w, "		<div id='edit'><input type='hidden' name='password' value=''><input type='submit' value='Edit &raquo;'></div></form>")
				} else if public_view == false || strings.TrimSpace(r.FormValue("password")) != "" {
					fmt.Fprintf(w, "		<div id='edit'><input type='hidden' name='password' value='%s'><input type='submit' value='Edit &raquo;'></div></form>",strings.TrimSpace(r.FormValue("password")))
				} else {
					fmt.Fprintf(w, "		<div id='edit'><input type=\"password\" name=\"password\"><input type='submit' value='Edit &raquo;'></div></form>")
				}
			} else {
				fmt.Fprintf(w, "	<div id='content'><center><div id=\"error\"><div class=\"huge\">[</div><div id='error_container'><div class=\"big\">This Page is not public!</div><br>You think you are the owner?<br><br><form action='/view/%s' ><input type=\"password\" value=\"Password\" onFocus=\"if(this.value == 'Password'){ this.value = ''; }\" name=\"password\" id=\"password\"><input type=\"submit\" value='&raquo;'></form></div><div class=\"huge\">]</div></div></center></div>\n",p.Title)
			}

			fmt.Fprintf(w, "</body>\n</html>\n")
		} else {
			http.Redirect(w, r, "/view/Error", http.StatusFound)
		}
}

func renderEdit(w http.ResponseWriter ,p *Page, r *http.Request, markdown bool) {
		// Get Config from Body
		var boder = bytes.Split(p.Body,[]byte("~![META]!~"))
		// Replace Body with showable body
		if len(boder) > 1 {
			p.Body = boder[1]
			// Now lets look at the config meta
			var confer = bytes.Split(boder[0],[]byte("::"))
			// Init some trash vars
			var public_edit = false
			var public_edit_ = ""
			var public_view_ = ""
			// And korrigate them
			if bytes.Equal(confer[1],[]byte("true")) {
				public_edit = true
				public_edit_ = "checked"
			}
			if bytes.Equal(confer[0],[]byte("true")) {
				public_view_ = "checked"
			}
			// Get Password
			var password = strings.TrimSpace(string(confer[2]))

			// Print Default Body [HEAD PART]
			fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \n \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n" +
					"<html xmlns=\"http://www.w3.org/1999/xhtml\">\n" +
					"<head>\n" +
					"	<title>%s</title>\n" +
					"	<style type=\"text/css\">%s</style> \n " +
					" <script type='text/javascript'>function insertTab(a,b){var c=b.keyCode?b.keyCode:b.charCode?b.charCode:b.which;if(c==9&&!b.shiftKey&&!b.ctrlKey&&!b.altKey){var d=a.scrollTop;if(a.setSelectionRange){var e=a.selectionStart;var f=a.selectionEnd;a.value=a.value.substring(0,e)+\"\t\"+a.value.substr(f);a.setSelectionRange(e+1,e+1);a.focus()}else if(a.createTextRange){document.selection.createRange().text=\"\t\";b.returnValue=false}a.scrollTop=d;if(b.preventDefault){b.preventDefault()}return false}return true}</script>" +
					"</head>\n" +
					"<body>\n	<div id='header'>\n		<div class=\"right\"><a href='/'>%s</a></div>\n		<h1>Edit: <a href=\"/view/%s\">%s</a></h1>\n	</div>\n", p.Title+" ~ "+wikiname, css_, wikiname, p.Title, p.Title)

			if public_edit || strings.TrimSpace(r.FormValue("password")) == password {
				fmt.Fprintf(w, "	<form action='/save/%s' >\n		<div id='edit_text'><textarea onkeydown=\"insertTab(this, event);\" name='body'>%s</textarea></div>\n", p.Title, p.Body)
				fmt.Fprintf(w, "	\n		<div id='edit_right'>")
				fmt.Fprintf(w, "			<div id='edit'>\n				<input type='submit' class='long green' value='SAVE &raquo;'>\n<br/><br/>")
				fmt.Fprintf(w, "				<input type=\"checkbox\" name='public_view' %s> Public view?<br/>",public_view_)
				fmt.Fprintf(w, "				<input type=\"checkbox\" name='public_edit' %s> Public edit?<br/><br/>",public_edit_)
				fmt.Fprintf(w, "				<input type=\"hidden\" name='password_' value='%s'>",strings.TrimSpace(r.FormValue("password")))
				fmt.Fprintf(w, "				<fieldset><legend>New Password</legend><input type=\"password\" value='%s' name='new_password'></fieldset>",strings.TrimSpace(r.FormValue("password")))
				fmt.Fprintf(w, "				</form><br/><form action='/remo/%s'><input type='hidden' name='pass' value='%s'><input type='submit' OnClick='JavaScript:if(!confirm(\"Really Remove this PagE?\")){ return false; }' clasS='long' value='Remove'></form>",p.Title,strings.TrimSpace(r.FormValue("password")))
				fmt.Fprintf(w, "			</div>\n	</div>")

				fmt.Fprintf(w, "</body>\n</html>\n")
			} else {
				fmt.Fprintf(w, "	<div id='content'><center><div id=\"error\" class='.wide'><div id='error_container'><div class=\"big\">And what do you think you do?<br /><br /><big>O.o</big></div><br><small>Maybe it was just a wrong password... [<a href=\"/view/%s\">back</a>]</small></div></center></div>\n",p.Title)
			}
	} else {
		http.Redirect(w, r, "/view/Error", http.StatusFound)
	}


}

func renderCreate(w http.ResponseWriter, r *http.Request, title string) {
		// Get Config from Body

		if title != "CreatePage" {
			fmt.Fprintf(w, "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \n \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n" +
					"<html xmlns=\"http://www.w3.org/1999/xhtml\">\n" +
					"<head>\n" +
					"	<title>New Wiki/Page ~ %s</title>\n" +
					"	<style type=\"text/css\">%s</style> \n " +
					"</head>\n" +
					"<body>\n	<div id='header'>\n		<div class=\"right\"><a href='/'>%s</a></div>\n		<h1>New Wiki/Page</h1>\n	</div>\n",wikiname,css_,wikiname)

			fmt.Fprintf(w, "<div id='content'><center><div id='new'>" +
				"<fieldset><form action='/cnew/CreatePage' method='get'>" + 
				"<legend>Title</legend><input type=\"text\" name='title' value=\"%s\"><br>" +
				"<legend>Password</legend><input name='pass' type=\"password\"><br><br>" +
				"Public view?<input type=\"checkbox\" name='view'>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" +
				"Public edit?<input type=\"checkbox\" name='edit'><br><br>" +
				"<input type=\"submit\" value=\"Create new Page!\"><br>" +
				"</form></fieldset></div></center></div>\n",title)

			fmt.Fprintf(w, "</body>\n</html>\n")
		} else {
		  _, err := loadPage(title)
		  if err != nil {
				public_view := "false"
				public_edit := "false"
				if r.FormValue("view") == "on" {
					public_view = "true"
				}
				if r.FormValue("edit") == "on" {
					public_edit = "true"
				}
			 	new_password := r.FormValue("pass")
				body := public_view+"::"+public_edit+"::"+new_password+"\n~![META]!~\n"+string(default_)
				p := &Page{Title: r.FormValue("title"), Body: []byte(body)}
				err := p.save()
				if err != nil {
				    http.Error(w, err.Error(), http.StatusInternalServerError)
				    return
				}
				http.Redirect(w, r, "/view/"+r.FormValue("title")+"?password="+new_password, http.StatusFound)
		} else {
				http.Redirect(w, r, "/view/"+r.FormValue("title"), http.StatusFound)
		}
	}
}


func getTitle(w http.ResponseWriter, r *http.Request) (title string, err error) {
    title = r.URL.Path[lenPath:]
    if !titleValidator.MatchString(title) {
        http.NotFound(w, r)
        err = errors.New("Invalid Page Title")
    }
    return
}



/*****************
** PAGE HANDLER **
*****************/

func makeMHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[len("/"):]
        fn(w, r, title)
    }
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        title := r.URL.Path[lenPath:]
        if !titleValidator.MatchString(title) {
						http.Redirect(w, r, "/", http.StatusFound)
            //http.NotFound(w, r)
            return
        }
        fn(w, r, title)
    }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/cnew/"+title, http.StatusFound)
        return
    }
    renderView(w, p, r, true)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/cnew/"+title, http.StatusFound)
        return
    }
    renderEdit(w, p, r, false)
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
		page, err := loadPage(title)
		if err != nil || title == "" {
    	renderCreate(w, r, title)
			return
		}
		http.Redirect(w, r, "/view/"+page.Title, http.StatusFound)
}

func cnewHandler(w http.ResponseWriter, r *http.Request, title string) {
		page, err := loadPage(title)
		if err != nil {
    	renderCreate(w, r, title)
			return
		}
		http.Redirect(w, r, "/view/"+page.Title, http.StatusFound)
}

func remoHandler(w http.ResponseWriter, r *http.Request, title string) {
		page, _ := loadPage(title)
		boder := bytes.Split(page.Body,[]byte("~![META]!~"))
		confer := bytes.Split(boder[0],[]byte("::"))
		if strings.TrimSpace(r.FormValue("pass")) == strings.TrimSpace(string(confer[2])) {
		  os.Remove("pages/"+title+".txt")
		  http.Redirect(w, r, "/?saved=succes", http.StatusFound)
		} else {
		  http.Redirect(w, r, "/view/"+title+"?password="+strings.TrimSpace(r.FormValue("password_"))+"&saved=fail", http.StatusFound)
		}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
		page, _ := loadPage(title)
		boder := bytes.Split(page.Body,[]byte("~![META]!~"))
		confer := bytes.Split(boder[0],[]byte("::"))
		if strings.TrimSpace(r.FormValue("password_")) == strings.TrimSpace(string(confer[2])) {
		  body := r.FormValue("body")
			public_edit := "false"
			public_view := "false"
		  if r.FormValue("public_view") == "on" {
				public_view = "true"
			}
		  if r.FormValue("public_edit") == "on" {
				public_edit = "true"
			}
		 	new_password := r.FormValue("new_password")
			body = public_view+"::"+public_edit+"::"+new_password+"\n~![META]!~\n"+body
		  p := &Page{Title: title, Body: []byte(body)}
		  err := p.save()
		  if err != nil {
		      http.Error(w, err.Error(), http.StatusInternalServerError)
		      return
		  }
		  http.Redirect(w, r, "/view/"+title+"?password="+new_password+"&saved=succes", http.StatusFound)
		} else {
		  http.Redirect(w, r, "/view/"+title+"?password="+strings.TrimSpace(r.FormValue("password_"))+"&saved=fail", http.StatusFound)
		}
}

/*****************
**     MAIN     **
*****************/

func main() {
		// Load Default Libs
		fmt.Println("~~~~~~~[ "+wikiname+" ]~~~~~~~")
		fmt.Println(":: VERSION = 0.1")
		fmt.Println(":: AUTHOR = <Roman P> setamagiga@gmail.com")
		fmt.Print(":: Loading CSS")
		css_ = loadLib("style.css")
		fmt.Println("... DONE")
		fmt.Print(":: Loading Default Page")
		default_ = loadLib("default.txt")
		fmt.Println("... DONE")

    http.HandleFunc("/", makeMHandler(homeHandler))
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    http.HandleFunc("/cnew/", makeHandler(cnewHandler))
    http.HandleFunc("/remo/", makeHandler(remoHandler))

    http.ListenAndServe(":8080", nil)
}

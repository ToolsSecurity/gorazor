# GoRazor

[![Build Status](https://travis-ci.org/sipin/gorazor.svg?branch=master)](https://travis-ci.org/sipin/gorazor)

GoRazor is the Go port of the razor view engine originated from [asp.net in 2011](http://weblogs.asp.net/scottgu/archive/2010/07/02/introducing-razor.aspx). In summary, GoRazor's:

* Concise syntax, no delimiter like `<?`, `<%`, or `{{`.
  * Original [Razor Syntax](http://www.asp.net/web-pages/tutorials/basics/2-introduction-to-asp-net-web-programming-using-the-razor-syntax) & [quick reference](http://haacked.com/archive/2011/01/06/razor-syntax-quick-reference.aspx/) for asp.net.
* Able to mix go code in view template
  * Insert code block to import & call arbitrary go modules & functions
  * Flow controls are just Go, no need to learn another mini-language
* Code generation approach
  * No reflection overhead
  * Go compiler validation for free
* Strong type view model
* Embedding templates support
* Layout/Section support

# Usage

Install:

```sh
go get github.com/sipin/gorazor
```

Usage:

`gorazor template_folder output_folder` or
`gorazor template_file output_file`

# Syntax

## Variable

* `@variable` to insert **string** variable into html template
  * variable could be wrapped by arbitrary go functions
  * variable inserted will be automatically [esacped](http://golang.org/pkg/html/template/#HTMLEscapeString)

```html
<div>Hello @user.Name</div>
```

```html
<div>Hello @strings.ToUpper(req.CurrentUser.Name)</div>
```

Use `raw` to skip escaping:

```html
<div>@raw(user.Name)</div>
```

Only use `raw` when you are 100% sure what you are doing, please always be aware of [XSS attack](http://en.wikipedia.org/wiki/Cross-site_scripting).

## Flow Control

```php
@if .... {
	....
}

@if .... {
	....
} else {
	....
}

@for .... {

}

@{
	switch .... {
	case ....:
	      <p>...</p>
	case 2:
	      <p>...</p>
	default:
	      <p>...</p>
	}
}
```

Please use [example](https://github.com/sipin/gorazor/blob/master/examples/tpl/home.gohtml) for reference.

## Code block

It's possible to insert arbitrary go code block in the template, like create new variable.

```html
@{
	username := u.Name
	if u.Email != "" {
		username += "(" + u.Email + ")"
	}
}
<div class="welcome">
<h4>Hello @username</h4>
</div>
```

It's recommendation to keep clean separation of code & view. Please consider move logic into your code before creating a code block in template.

## Declaration

The **first code block** in template is strictly for declaration:

* imports
* model type
* layout

like:

```go
@{
	import  (
		"kp/models"   //import `"kp/models"` package
		"tpl/layout/base"  //Use tpl/layout package's **base func** for layout
	)
	var user *models.User //1st template param
	var blog *models.Blog //2nd template param
}
...
```

**first code block** must be at the beginning of the template, i.e. before any html.

Any other codes inside the first code block will **be ignored**.

import must be wrapped in `()`, `import "package_name"` is not yet supported.

The variables declared in **first code block** will be the models of the template, i.e. the parameters of generated function.

If your template doesn't need any model input, then just leave it blank.

## Helper / Include other template

As gorazor compiles templates to go function, embedding another template is just calling the generated function, like any other go function.

However, if the template are designed to be embedded, it must be under `helper` namespace, i.e. put them in `helper` sub-folder.

So, using a helper template is similar to:

```html

@if msg != "" {
	<div>@helper.ShowMsg(msg)</div>
}

```

GoRazor won't HTML escape the output of `helper.XXX`.

Please use [example](https://github.com/sipin/gorazor/blob/master/examples/tpl/home.gohtml) for reference.

## Layout & Section

The syntax for declaring layout is a bit tricky, in the example mentioned above:

```go
@{
	import  (
		"tpl/layout/base"
	)
}
```

`"tpl/layout/base"` **is not** a package, it's actually referring to `"tpl/layout"` package's **Base** function, which should be generated by `tpl/layout/base.gohtml`.

GoRazor is using the second last part of namespace `layout` as a magic string to decide if the import is for layout declaration or normal import.

A layout file `tpl/layout/base.gohtml` may look like:

```html
@{
	var body string
	var sidebar string
	var footer string
	var title string
	var css string
	var js string
}

<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8" />
	<title>@title</title>
</head>
<body>
    <div class="container">@body</div>
    <div class="sidebar">@sidebar</div>
    <div class="footer">@footer</div>
	@js
  </body>
</html>
```

It's just a usual gorazor template, but:

* First param must be `var body string` (As it's always required, maybe we could remove it in future?)
* All params **must be** string, each param is considered as a **section**, the variable name is the **section name**.
* Under `layout` package, i.e. within "layout" folder.

A template using such layout `tpl/index.gohtml` may look like:

```html
@{
	import (
		"tpl/layout"
	)
}

@section footer {
	<div>Copyright 2014</div>
}

<h2>Welcome to homepage</h2>
```

With the page, the page content will be treated as the `body` section in layout.

The other section content need to be wrapped with
```
@section SectionName {
	....
}
```

The template doesn't need to specify all section defined in layout. If a section is not specified, it will be consider as `""`.

Thus, it's possible for the layout to define default section content in such manner:

```html
@{
	var body string
	var sidebar string
}

<body>
    <div class="container">@body</div>
    @if sidebar == "" {
    <div class="sidebar">I'm the default side bar</div>
	} else {
    <div class="sidebar">@sidebar</div>
	}
</body>
```

* A layout should be able to use another layout, it's just function call.

# Conventions

* Template **folder name** will be used as **package name** in generated code
* Template file name must has the extension name `.gohtml`
* Template strip of `.gohtml` extension name will be used as the **function name** in generated code, with **first letter Capitalized**.
  * So that the function will be accessible to other modules. (I hate GO about this.)
* Helper templates **must** has the package name **helper**, i.e. in `helper` folder.
* Layout templates **must** has the package name **layout**, i.e. in `layout` folder.

# Example

Here is a simple example of [gorazor templates](https://github.com/sipin/gorazor/tree/master/examples/tpl) and the corresponding [generated codes](https://github.com/sipin/gorazor/tree/master/examples/gen).

# FAQ

## IDE / Editor support?

* Sublime Text 2/3 Syntax Highlight: search & install `GoRazor` via Package Control.

## How to auto re-generate when gohtml file changes?

We may add `gorazor watch` cmd after Go 1.3 which has official [fsnotify](https://docs.google.com/document/d/1xl_aRcCbksFRmCKtoyRQG9L7j6DIdMZtrkFAoi5EXaA/edit) support.

Currently, we are using below scripts to handle this issue on mac:

* gorazor_watch.sh
```bash
#!/bin/bash

gorazor tpl src/tpl
watchmedo shell-command --patterns="*.gohtml" --recursive --command='python gorazor.py ${watch_src_path}'
```

* gorazor.py
```python
import sys, os

path = sys.argv[1]
os.system("gorazor " + path + " " + path.replace("/tpl/", "/src/tpl/")[:-4])
```

* [watchmedo](https://github.com/gorakhargosh/watchdog)

# Credits

The very [first version](https://github.com/sipin/gorazor/releases/tag/vash) of GoRazor is essentially a hack of razor's port in javascript: [vash](https://github.com/kirbysayshi/vash), thus requires node's to run.

GoRazor has been though several rounds of refactoring and it has completely rewritten in pure Go. Nonetheless, THANK YOU [@kirbysayshi](https://github.com/kirbysayshi) for Vash! Without Vash, GoRazor may never start.

# Todo

* Add tools, like monitor template changes and auto re-generate
* Add default html widgets
* Add more usage examples
* Generate more function overloads, like accept additional buffer parameter for write

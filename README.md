# go-weather-widget
Eccosia's weather widget wingding challenge for WWG Berlin

## To run the application (we start the challenge with a bunch of failing tests):
```
go test ./... && go build . && ./go-weather-widget -api_key=YOUR_KEY_TO_WORLD_WEATHER_ONLINE
```
### Steps to solve the challenge:

#### Build the layout templates structure from outside in:

![alt text](https://github.com/wwgberlin/go-weather-widget/blob/master/layout.jpg "layout")


At the very top our application always renders a template named "layout" - which we define in the tmpl/templates/layouts/layout.tmpl file.

Layout:

1. Start the template by defining the "layout" template `{{define "layout"}} {{end}}`
2. Inside layout render the enclosing <HTML></HTML> tags.
3. Inside the <HTML></HTML> element render the "head" template (defined in the head.tmpl) - `{{template "head" .}}` - The dot is to pass the arguments to the head template.
4. Add the wrapping <BODY></BODY> element and inside the <HTML> element. 
5. Inside your body element, render the template "content" (don't forget to pass in the arguments).
6. Define empty "content" and "head" at the very end of your file (to prevent errors and allow rendering with default definitions of those templates.

Now we are ready to implement the "head" template in the tmpl/templates/layouts/head.tmpl file.

Head:
1. Define the "head" template in head.tmpl.
2. Inside the head template render the <HEAD></HEAD> wrapping HTML tags
3. Inside the HEAD element render the "title" template.
4. Inside the HEAD element render the "styles" template.
5. Add empty default "styles" and "title" templates at the end of the file.
	
#### Now let's look at our weather widget at templates/widget.tmpl
There's a comment inside where your implementation should go.
Your task here is to register the function clothes in the helpers FuncMap and call it with the descriptiona and the celsius values. Range over the returned values and for each render a <div/> element with each of the clothings.

#### Now let's imeplement our renderer (tpl/renderer.go)
1. Implement BuildTemplate according to instructions
2. Implement RenderTemplate according to instructions

#### Now let's imeplement our widgetHandler (handlers.go)
Follow the instructions in the file to implement the widgetHandler

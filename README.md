# rightempl
Go templating done right.

text/template is slow & not composable.
QuickTemplate requires a generate step which means no live editing. 

I took https://github.com/valyala/fasttemplate for fast edits & went further:

__Templates loaded, but auto-reload__
   Instead of reading each time, lets keep them in memory.

___Composition___
   Really go crazy with functions which render their own template. Property passing would be important.

___Frame-always___
   If you always wrap your page with the same head, then just send it before any computation. That'll get the browser heading this thing out the door. 
   
___Scriptable caching___ 
   Renders are often worth keeping around for a while, so why make the template system opaque to that? Allow an interceptor. 

___Useful Functions___
   Local statics can be served with a hash if we know where the static folders live. 

# rightmpl
Go templating done right.

text/template is slow & not composable.
QuickTemplate requires a generate step which means no live editing. 

https://github.com/valyala/fasttemplate offers fast edits, but lets go further:

_Today_

__Templates loaded, but auto-reload__
   Instead of reading each time, lets keep them in memory.

__Variable Substitution works like fasttemplate allows_

_Coming Soon_

___New with funcs___
If all functions are known before parsing, then dynamic dispatch stacks can form.

Plan:
#. Get subtemplate working as a function
#. Get boolean expressions working
#. "IF" functions eval a boolean expression and run a subtemplate
#. reflect-less types?   .I.count   .F.floatexample  .T.text  .B.boolexample .S.Sliceexample   .M.Mapexample
        .T.add("key", "value") <-- nah, still needs a map at runtime.
           or just make it easy with map[string]interface{}  and do work in post?

___Composition___
   Really go crazy with functions which render their own template. Property passing would be important.

___Frame-always___
   If you always wrap your page with the same head, then just send it before any computation. That'll get the browser heading this thing out the door. 
   
___Useful Functions___
   Local statics can be served with a hash if we know where the static folders live. 

_Possible_

___Full text/template Compat__ It should be possible to reuse the parser

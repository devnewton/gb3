{
var backend2html = {};
backend2html.tagStack = [];

backend2html.encode = function encode(e) {
    return e.replace(/[.]/g, function(e) {
        return `&#${e.charCodeAt(0)};`;
    });
};

}

post
 = items:postItem*
 {
 	var result = items.join("");
    while(backend2html.tagStack.length > 0) {
    	result += `</${backend2html.tagStack.pop()}>`;
    }
    return result;
 }
 
postItem
 = url
 / totoz
 / bigorno
 / norloge
 / openTag
 / closeTag
 / xmlEntities
 / xmlSpecialChar
 / .

 xmlEntities
  = $("&" ("lt" / "gt" / "amp" ) ";")

xmlSpecialChar
 = (lt / gt / amp / quot / apos)
 
lt
 = "<"
 { return "&lt;"; }
 
gt
 = ">"
 { return "&gt;"; }
 
amp
 = "&"
 { return "&amp;"; }

quot
 = '"'
 { return "&quot;"; }
 
apos
 = "'"
 { return "&apos;"; }
 
url
 = "<a href=\""? urlStr:$((("http" "s"?) / "ftp") "://" ([^< \t\r\n\"])+) ("\">" [^<]+ "</a>")?
 { 
   let url = new URL(urlStr)
   return `<a href="${encodeURI(url)}" target="_blank" rel="noreferrer">${url.hostname}</a>`;}

openTag
 = "<" tag:validFormatTag ">"
 {
 	backend2html.tagStack.push(tag);
 	return `<${tag}>`;
 }

closeTag
 = "</" tag:validFormatTag ">"
 {
 	var result = "";
 	for(;;) {
      var poppedTag = backend2html.tagStack.pop();
      if(poppedTag == undefined) {
      	break;
      }
      if( poppedTag == tag) {
      	result += `</${tag}>`;
        break;
      } else {
      	result += `</${poppedTag}>`;
      }
    }
    return result;
 }

validFormatTag
 = (spoiler / "b" / "i" / "s" / "u" / "tt" / "code" )

spoiler
 = "spoiler"
 { return "gc2-spoiler"; }

invalidOpenTag
 = "<" tag:invalidTag ">"
 { return `&lt;${tag}&gt;`; }

invalidCloseTag
 = "</" tag:invalidTag ">"
 { return  `&lt;/${tag}&gt;`; }
 
invalidTag
 = [A-Za-z] (xmlSpecialChar / [^>])*
 
tagAttributes
 = attributes:(separator:" " attribute:tagAttribute { return attribute;})*
 {  var result = {};
 	for(var a in attributes) {
    	result[attributes[a].name] = attributes[a].value;
    }
 	return result;
 }

tagAttribute
 = name:$[a-z]+ value:("=\"" value:$[^"]* "\"" {return value;} )?
 {return { name: name, value: value}}

bouchot
 = $("@" [a-z]+)


norloge
 = fullNorloge / longNorloge / normalNorloge / shortNorloge

fullNorloge
 = y: norlogeYear "-" m: norlogeMonth "-" d:norlogeDay [T# ] h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds b:bouchot?
 {
 let time = `${h}:${mi}:${s}`;
 let bouchot = b && ` bouchot="${b.substr(1)}"` || "";
 return `<gc2-norloge title="${y}-${m}-${d} ${time}"${bouchot}>${time}${b || ""}</gc2-norloge>`;
 }

 longNorloge
 = m: norlogeMonth "/" d:norlogeDay "#" h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds b:bouchot?
 {
 let time = `${h}:${mi}:${s}`;
 let bouchot = b && ` bouchot="${b.substr(1)}"` || "";
 return `<gc2-norloge title="${m}-${d} ${time}"${bouchot}>${time}${b || ""}</gc2-norloge>`;
 }
 
norlogeYear
 = digits: [0-9]+
 { return digits.join(""); }
 
norlogeMonth
 = first: [0-1] last: [0-9]
 { return first + last; }

norlogeDay
 = first: [0-3] last: [0-9]
 { return first + last; }

normalNorloge
 = h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds b:bouchot?
 {
 let time = `${h}:${mi}:${s}`;
 let bouchot = b && ` bouchot="${b.substr(1)}"` || "";
 return `<gc2-norloge title="${time}"${bouchot}>${time}${b || ""}</gc2-norloge>`;
 }
 
shortNorloge
 = h:norlogeHours ":" mi:norlogeMinutes b:bouchot?
 {
 let time = `${h}:${mi}`;
 let bouchot = b && ` bouchot="${b.substr(1)}"` || "";
 return `<gc2-norloge title="${time}"${bouchot}>${time}${b || ""}</gc2-norloge>`;
 }

norlogeHours
 = first: [0-2] last: [0-9]
 { return first + last; }
 
norlogeMinutes
 = first: [0-5] last: [0-9]
 { return first + last; }
 
norlogeSeconds
 = first: [0-5] last: [0-9]
 { return first + last; }

bigorno
 = spaces:$(inputStart / whitespaces) s2:whitespaces? bigorno:$[a-zA-Z0-9-_]+ ("<" / "&lt;") &(whitespaces / [<;:[,] / !.)
 { return `${spaces}<gc2-bigorno>${bigorno}</gc2-bigorno>`;}

totoz
  = first:"[:" totoz:$[^\]]+ third:"]"
  { return `<gc2-totoz>${backend2html.encode(totoz)}<img src="https://nsfw.totoz.eu/img/${encodeURI(totoz)}"></gc2-totoz>`; }
  
whitespaces
 = [ \t\r\n]

inputStart
 = & { return location().start.offset == 0; }

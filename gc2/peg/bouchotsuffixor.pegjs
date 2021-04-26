{
  let oldBouchotSuffix = options.oldBouchot ? `@${options.oldBouchot}` : "";
  let newBouchotSuffix = options.newBouchot ? `@${options.newBouchot}` : "";
  function suffixBouchot(b) {
    if(!b) {
    	return oldBouchotSuffix;
    }
    if(b === newBouchotSuffix) {
    	return "";
    }
    return b;
  }
}

post
 = items:postItem*
  {
 	return items.join("");
  }

postItem
 = url
 / totoz
 / norloge
 / .

url
 = $((("http" "s"?) / "ftp") "://" ([^< \t\r\n\"])+)

norloge
 = fullNorloge / longNorloge / normalNorloge / shortNorloge

fullNorloge
 = y: norlogeYear "-" m: norlogeMonth "-" d:norlogeDay sep:[T# ] h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds b:bouchot?
 {
 return `${y}-${m}-${d}${sep}${h}:${mi}:${s}${suffixBouchot(b)}`;
 }

 longNorloge
 = m: norlogeMonth "/" d:norlogeDay "#" h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds b:bouchot?
 {
 return `${m}/${d}#${h}:${mi}:${s}${suffixBouchot(b)}`;
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
   return `${h}:${mi}:${s}${suffixBouchot(b)}`;
 }
 
shortNorloge
 = h:norlogeHours ":" mi:norlogeMinutes b:bouchot?
 {
   return `${h}:${mi}${suffixBouchot(b)}`;

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

bouchot
 = $("@" [a-z]+)

totoz
  = $("[:" $[^\]]+ "]")

whitespaces
 = [ \t\r\n]
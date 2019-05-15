post
 = items:postItem*
  {
 	return items.join("");
  }

postItem
 = url
 / totoz
 / unbouchotedNorloge
 / .

url
 = $((("http" "s"?) / "ftp") "://" ([^< \t\r\n\"])+)

unbouchotedNorloge
 = n:norloge &(whitespaces / !.)
 {
   return `${n}@${options.bouchot}`
 }

norloge
 = fullNorloge / longNorloge / normalNorloge / shortNorloge

fullNorloge
 = y: norlogeYear "-" m: norlogeMonth "-" d:norlogeDay sep:[T# ] h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds
 {
 return `${y}-${m}-${d}${sep}${h}:${mi}:${s}`;
 }

 longNorloge
 = m: norlogeMonth "/" d:norlogeDay "#" h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds
 {
 return `${m}/${d}#${h}:${mi}:${s}`;
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
 = h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds
 {
   return `${h}:${mi}:${s}`;
 }
 
shortNorloge
 = h:norlogeHours ":" mi:norlogeMinutes
 {
   return `${h}:${mi}`;

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


totoz
  = $("[:" $[^\]]+ "]")

whitespaces
 = [ \t\r\n]
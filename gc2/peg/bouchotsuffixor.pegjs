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
 = y: norlogeYear "-" m: norlogeMonth "-" d:norlogeDay sep:[T# ] h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds
 {
 return `${y}-${m}-${d}${sep}${h}:${mi}:${s}@${options.bouchot}`;
 }

 longNorloge
 = m: norlogeMonth "/" d:norlogeDay "#" h:norlogeHours ":" mi:norlogeMinutes ":" s:norlogeSeconds
 {
 return `${m}/${d}#${h}:${mi}:${s}@${options.bouchot}`;
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
   return `${h}:${mi}:${s}@${options.bouchot}`;
 }
 
shortNorloge
 = h:norlogeHours ":" mi:norlogeMinutes
 {
   return `${h}:${mi}@${options.bouchot}`;

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

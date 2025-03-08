### cool title bro
Here is a sample blog post. This could be something like a long intro paragraph that has some intro stuff. like what is it all about really? Let me tell you cuz I know.

##### this would be a header
So here is the first paragraph and I'll get into some juicy deails right off the bat. Take a look at this stuff below here. Here is a code snippet:
```go
package helpers

import (
	"net/http"
)

func GetIpAddr(r *http.Request) string {
	// default to using custom set header
	ip := r.Header.Get("AK-First-External-IP")
	if ip != "" {
		return ip
	}

	// if above fails for some reason, use RemnoteAddr
	// note, that will just be the ip of the nginx reverse proxy
	ip = r.RemoteAddr
	return ip
}
```

##### another one
And now here's some more cool stuff to wet your eyeballs. Here is a sample [image](/home/myshkins/pics/synced_pics/selfie.jpeg), showing me.

Here is some **bold text** and some __italic text__. nice.

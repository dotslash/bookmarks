This is a simple chrome extension I use to bookmark urls and give a custom handle,
E.g. [bm.suram.in/r/latency-numbers](http://bm.suram.in/r/latency-numbers) redirects 
to [https://gist.github.com/jboner/2841832](https://gist.github.com/jboner/2841832)

The extension adds a popup next to the address bar, which on clicking lets the 
user (me) add custom handles to current tab url.If there are any handles to 
current tab url, the extension shows them all. 

Check the screenshots to see how the extention looks when
- Wrong secret entered when adding a handle
- A handle is successfully added
- There are already 2 custom handles

<img src='https://yesteapea.com/public/images/bm/wrong-secret.png' height=250>
<img src='https://yesteapea.com/public/images/bm/added-bm.png' height=250>
<img src='https://yesteapea.com/public/images/bm/two-bm.png' height=250>

## Installation
If you wish to make the plugin work with your own hosting of the bookmarks application,
change the `server_address` field in `popup.js`.

To install the app go to [`chrome://extensions/`](chrome://extensions/), turn on developer mode and `Load unpacked`


# FastgateRestarter
A quick and simple bot to restart the Fastweb modem/router Fastgate.
I preferred to use the flat model for this application because it contains a few lines of code

### **Analysis**
During my analysis I found that the web interface configuration (FRONT END) are written in Angular.
I analyzed the requests sent and found that:
 - **_** *[in the url]* contains the timestamp of the request
 - **password** *[in the url]* are base64 encoded
 - every page requires a **specific token**
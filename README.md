# WSproxy
Websockets to TCP sockets proxy

Websockets server that proxies traffic to a TCP server 

## Options
 -a Listen address (default "localhost")  
 -p Listen port (default 8080)  
 -r Remote TCP server address  
 -P Remote TCP server Port (default 21000)  
 -t Text encoding (default true)   
 -h Show help  

## Notes

The WebSocket protocol distinguishes between text and binary data messages. Text messages are interpreted as UTF-8 encoded text. The interpretation of binary messages is left to the application.
You can specify if you want to send text or binary data to Websockets connection using the **-t** flag.
By default the messages are sent in text (UTF-8) mode

## Example

wsproxy -a localhost -p 8080 -r myserver.com -P 33000 -t true

Wsproxy starts to listening connections in localhost:8080 and redirect traffic to myserver.com:33000 with UTF-8 encoding

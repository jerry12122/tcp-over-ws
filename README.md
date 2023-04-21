# tcp-over-ws
實現一個Server-Client架構，主要功能為tcp-over-websocket代理  
連接為 A服務-Client-Server-B服務  
Client-Server中間是websocket連接  
A服務-Client、Server-B服務則為tcp連接  

Client
* 上面有一個tcp listener(tcp://127.0.0.1:3333)持續監聽
* 只有tcp listener建立tcp Connection 的事件發生時，才建立WebSocket Connection (ws://ws.local:8083/)
* 當WebSocket Connection建立之後，將tcp Connection 傳入的訊息Read出來並Write到WebSocket Connection
* tcp Connection 不會主動斷開
* 當tcp Connection 斷開後10秒沒有再建立 tcp Connection 就斷開ws Connection 

Server
* 包括一個WebSocket Listener(ws://127.0.0.1:8082)持續監聽
* Server作為WebSocket的監聽端，是只有在觸發接收事件的時候取得連接，不會主動建立連接
* WebSocket Listener收到傳入訊息的時候，建立tcp Connection(tcp://127.0.0.1:3333)並把訊息write到tcp Connection



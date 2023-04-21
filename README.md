# tcp-over-ws
## 簡介
> 實現一個Server-Client架構，主要功能為TCP-over-websocket代理
> 連接為 A服務-Client-Server-B服務
> Client-Server中間是websocket連接
>　A服務-Client、Server-B服務則為TCP連接

## Client
* 上面有一個TCP Listener(TCP://127.0.0.1:3333)持續監聽
* 當TCP Listener有新的TCP Connection，才會建立新的WebSocket Connection(ws://ws.local:8082/)
* WebSocket Connection會持續連接到TCP Connection被斷開
* 將TCP Connection 傳入的訊息Read出來輸出(log.Println)並Write到WebSocket Connection
* TCP Connection WebSocket Connection 都不會主動斷開
* 當TCP Connection 被斷開後10秒沒有再建立 TCP Connection 才會斷開 WebSocket Connection
* WebSocket Connection斷開後，直到有新的TCP Connection才會重新建立WebSocket Connection

## Server
* 上面有一個WebSocket Listener(ws://127.0.0.1:8082)持續監聽
* 偵測WebSocket Listener建立WebSocket Connection後，才建立新的TCP Connection(TCP://127.0.0.1:3333)
* 建立的WebSocket Connection收到傳入訊息的時候，通過剛剛建立的TCP Connection並把訊息write到TCP Connection
* 保留建立的WebSocket Connection，若TCP Connection有回傳訊息時，通過保留的WebSocket Connection回傳
* 若TCP Connection沒有回傳訊息繼續轉傳WebSocket Connection收到傳入訊息
* TCP Connection WebSocket Connection 都不會主動斷開




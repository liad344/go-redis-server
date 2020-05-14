# Go-redis-server

Partial redis implementation using go! 

redis benchmark- get/set
- 26,500/sec on intel 4600U (Windows, Up from 16000~ on first benchmark )
- 72,000/sec on intel Xeon(R) Platinum 8171M (Linux)  

## TODO 
  - fix reading from client
  - add unsafe get/set
  - add snapshot? 
  - add config file

# Go-redis-server

Partial redis implementation using go! 

redis benchmark- get/set
- 28,000/sec on intel 4600U (Windows, Up from 16000~ on first benchmark, 26500 with 2-4 goroutines)
- 71,000/sec on intel 10210U (4C/8T with one gorutine , 68,000 with 4-8 goruotines)
## TODO 
  - fix reading from client
  - add unsafe get/set
  - add snapshot? 
  - add config file

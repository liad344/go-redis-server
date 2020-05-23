# Go-redis-server

Partial redis implementation using go! 

redis benchmark- get/set
- 28,000/sec on intel 4600U (2C/4T, localhost , first bench 16,000~)
- 71,000/sec on intel 10210U (4C/8T, localhost)
- 110,000/sec on intel 8700 (6C/12T over 1gb network)
## TODO 
  - fix reading from client
  - add unsafe get/set
  - add snapshot? 
  - add config file

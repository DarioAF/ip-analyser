<div align="center">
    <h1>IP Analyser 🔍</h1>
    <sub>Simple project to analyse and return basic info about a given ip address.</sub>
</div>

---
## Set up:
```bash
docker-compose up --build
```
This command will startup the main app and a [redis](https://redis.io/) server in your machine.<br />
(Make sure you have [docker](https://www.docker.com/) properly installed)<br />
<br />By default, the app will initialize at `http://localhost:8080`<br />
By calling the root route, you should see the following response:

```json
{
    "status": "UP & Running"
}
```

---
## Endpoints:

### 👨‍💻 User
```scala
[POST] "/user"
```
Takes an object `User` with the only information we have at this point (the ip address) and outputs the analysis result, containing all the aditional data we have collected from that ip address.

#### Example
input object:
```json
{
    "ip": "83.44.196.93",
}
```

output: 
```json
{
    "ip": "83.44.196.93",
    "time": "20/06/2020 02:42:46",
    "country": "Spain",
    "iso_country": "ES",
    "distance": 10274,
    "is_aws": false
}
```
(We will talk about each of this items later on)

### 📈 Analytics
The app collect data about each of the analysed ip addresses and makes it available through the following services<br />

To simplify examples, imagine we collected this info:

| IP | Country | Distance | Count |
| :---: | :---: | :---: | :---: |
| 1.1.1.1  | Argentina  | 0  | 10 |
| 2.2.2.2  | Brazil  | 2821  | 100 |
| 3.3.3.3  | Spain  | 10274  | 50 |
| 4.4.4.4  | Spain  | 10274  | 30 |

```scala
[GET] "/nearest"
[GET] "/farthest"
```
Returns the nearest / farthest distance from **Argentina** to wherever the service was invoked in *km* <br />

#### Example

output for `/nearest` in our example:
```json
{
    "country": "AR",
    "distance": 0
}
```
output for `/farthest` in our example:
```json
{
    "country": "ES",
    "distance": 10274
}
```
**Note:** that if there are two or more countries with the same distance from AR, only the one with more requests will be returned

```scala
[GET] "/avg-requests/{country-code}"
```
Returns the request average for the specified country

#### Example

output for `/avg-requests/ES` in our example:
```json
{
    "avg": 40
}
```
```go
(30+50)/2 = 40
```
output for `/avg-requests/BR` in our example:
```json
{
    "avg": 100
}
```
```go
(100)/1 = 100
```

---
## Tests:
Run the included tests by running:
```
go test
```
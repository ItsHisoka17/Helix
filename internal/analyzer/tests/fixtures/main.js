const redis = require("ioredis");

let app = redis.createClient();

app.listen(80);
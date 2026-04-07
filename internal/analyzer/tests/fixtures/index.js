const express = require("express");

let app = express();

app.get("/", (req, res) => {
    let status = res.status ?? 200;
    res.json({
        message: "Hello from Helix",
        status
    })
})

app.listen(3000);
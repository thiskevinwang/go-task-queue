import express from "express";

const app = express();

app.get("/", (req, res) => {
  res.send(
    "This is running nodemon, and Docker Compose resyncs on file changes"
  );
});

app.listen(3000, () => {
  console.log("Server running on port 3000");
});

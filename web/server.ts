import { Hono } from "hono";
import { serve } from "@hono/node-server";

const app = new Hono();

app.get("/", (c) => {
  return c.text("Hello World");
});

serve(app, (info) => {
  console.log(`Server is running at ${info.address}:${info.port}`);
});

import { Hono } from "hono";
import { timing } from "hono/timing";
import { logger } from "hono/logger";

import { serve } from "@hono/node-server";

const app = new Hono();

// https://hono.dev/middleware/builtin/timing
app.use("*", timing());

//https://hono.dev/middleware/builtin/logger
app.use("*", logger());

app.get("/", (c) => {
  return c.text("Hello World");
});

serve(app, (info) => {
  console.log(`Server is running at ${info.address}:${info.port}`);
});

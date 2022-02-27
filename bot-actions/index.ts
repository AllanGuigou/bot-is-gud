import * as gcp from "@pulumi/gcp";

import * as express from "express";
import * as bodyParser from "body-parser";

import { conversation } from "@assistant/conversation";
import axios from "axios";

const { DISCORD_ROLE, DISCORD_WEBHOOK_URL } = process.env;

// https://actions-on-google.github.io/assistant-conversation-nodejs/3.7.0/interfaces/conversation_conv.conversationv3options.html
const _function = new gcp.cloudfunctions.HttpCallbackFunction(
  "at-discord-role-fulfillment",
  {
    runtime: "nodejs16",
    availableMemoryMb: 128,
    environmentVariables: {
      DISCORD_ROLE,
      DISCORD_WEBHOOK_URL,
    },
    callbackFactory: () => {
      const { DISCORD_ROLE, DISCORD_WEBHOOK_URL } = process.env;
      const app = conversation({ debug: true });

      app.handle("at_discord_role", async (conv) => {
        await axios.post(DISCORD_WEBHOOK_URL!, {
          content: DISCORD_ROLE!,
        });
        conv.add("Notified: " + DISCORD_ROLE);
      });

      const expressApp = express();

      expressApp.use(bodyParser.json());

      expressApp.post("/fulfillment", app);

      return app;
    },
  }
);

const invoker = new gcp.cloudfunctions.FunctionIamMember("invoker", {
  project: _function.function.project,
  region: _function.function.region,
  cloudFunction: _function.function.name,
  member: "allUsers",
  role: "roles/cloudfunctions.invoker",
});

export const url = _function.httpsTriggerUrl;

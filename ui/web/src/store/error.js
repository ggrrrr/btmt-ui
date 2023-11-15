import { defineStore } from "pinia";

import { useConfig } from "@/store/app";
const config = useConfig;

export const useErrorStore = defineStore({
  id: "error",
  state: () => ({
    type: "",
    message: "",
    errors: "",
  }),
  actions: {
    reset() {
      this.type = "";
      this.message = "";
      this.errors = {};
    },
    authError(msg) {
      this.type = "auth";
      this.message = msg;
      this.errors = {};
    },
    systemErr(msg, err) {
      this.type = "system";
      this.message = msg;
      this.errors = { system: err };
    },
    networkErr(msg, err) {
      this.type = "network";
      this.message = msg;
      this.errors = {
        baseUrl: config.BASE_URL,
        network: err,
      };
    },
    invalidResponse(message, error, response) {
      console.log("json.error", error);
      console.log(error);
      console.log(response);

      this.type = "system";
      this.message = message;
      this.errors = {
        error: error,
        response: `[${response.status}]: ${response.statusText}`,
      };
    },
  },
});

import { defineStore } from "pinia";

import { useConfig } from "@/store/app";
const config = useConfig;

export const useErrorStore = defineStore({
  id: "error",
  state: () => ({
    type: "",
    message: "",
    errors: "",
    show: false,
  }),
  actions: {
    alertType() {
      if (this.type == "network") {
        return "error";
      }
      if (this.type == "auth") {
        return "error";
      }
      return "warning";
    },
    Hide() {
      this.show = false;
    },
    clean() {
      this.type = "";
      this.message = "";
      this.errors = {};
      this.show = false;
    },
    authError(msg) {
      this.type = "auth";
      this.message = msg;
      this.errors = {};
      this.show = true;
    },
    systemErr(msg, err) {
      this.type = "system";
      this.message = msg;
      this.errors = { system: err };
      this.show = true;
    },
    networkErr(msg, err) {
      this.show = true;
      this.type = "network";
      this.message = msg;
      this.errors = {
        baseUrl: config.BASE_URL,
        network: err,
      };
    },
    invalidResponse(message, error, response) {
      this.show = true;
      console.log("json.error", error);
      console.log(error);
      console.log(response);

      this.alertType = "error";
      this.type = "system";
      this.message = message;
      this.errors = {
        error: error,
        response: `[${response.status}]: ${response.statusText}`,
      };
    },
  },
});

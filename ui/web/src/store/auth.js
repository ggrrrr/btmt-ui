// Utilities
import { defineStore } from "pinia";

import { useConfig } from "@/store/app";
import { useErrorStore } from "@/store/error";

const config = useConfig;

let errorStore = useErrorStore();

export const useLoginStore = defineStore({
  id: "auth",
  state: () => ({
    email: localStorage.getItem("email"),
    token: localStorage.getItem("token"),
    // error: "",
  }),
  actions: {
    createHeaders() {
      if (this.token) {
        const bearer = "Bearer " + this.token;
        return {
          "Content-Type": "application/json",
          authorization: bearer,
        };
      }
      return {
        "Content-Type": "application/json",
      };
    },
    logout() {
      this.token = "";
      this.email = "";
    },
    resetLogin() {
      this.token = "";
      this.email = "";
      localStorage.setItem("email", "");
      localStorage.setItem("token", "");
    },
    async validateRequest() {
      const url = config.BASE_URL + "/v1/auth/validate";
      const requestOptions = {
        method: "POST",
        headers: this.createHeaders(),
        body: JSON.stringify({}),
      };
      try {
        const response = await fetch(url, requestOptions);
        response
          .json()
          .then((data) => {
            console.log("data");
            console.log(data);
            if (data.code == 200) {
              errorStore.clean();
            } else if (data.code > 400 && data.code < 500) {
              this.resetLogin();
              errorStore.authError(data.message);
            } else {
              errorStore.networkErr(data.message, "");
            }
          })
          .catch((error) => {
            errorStore.invalidResponse("invalid response", error, response);
          });
      } catch (error) {
        console.log("error");
        console.log(error);
        errorStore.networkErr("unable to make request", error);
      }
    },
    async loginRequest(email, passwd) {
      const url = config.BASE_URL + "/v1/auth/login/passwd";
      // const url = config.BASE_URL + "/v1/nojson";
      const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        headers: this.createHeaders(),
        body: JSON.stringify({
          email: email,
          password: passwd,
        }),
      };
      try {
        const response = await fetch(url, requestOptions);
        response
          .json()
          .then((data) => {
            console.log("data");
            console.log(data);
            if (data.code == 200) {
              console.log("data === 200");
              console.log(data);
              this.email = data.payload.email;
              this.token = data.payload.token;
              localStorage.setItem("email", data.payload.email);
              localStorage.setItem("token", data.payload.token);
              errorStore.clean();
            } else if (data.code > 400 && data.code < 500) {
              console.log("code > 400 < 500");
              this.resetLogin();
              errorStore.authError(data.message);
            } else {
              console.log("code", data.code);
              errorStore.networkErr(data.message);
            }
          })
          .catch((error) => {
            errorStore.invalidResponse("invalid response", error, response);
          });
      } catch (error) {
        console.log("error");
        console.log(error);
        errorStore.networkErr("unable to make request", error);
      }
    },
  },
});

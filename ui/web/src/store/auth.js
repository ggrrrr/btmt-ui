// Utilities
import { defineStore } from "pinia";

import { useConfig } from "@/store/app";
import { useErrorStore } from "@/store/error";

const config = useConfig;

let errorStore = useErrorStore();

export const useLoginStore = defineStore({
  id: "auth",
  state: () => ({
    email: "",
    token: "",
    error: "",
  }),
  actions: {
    logout() {
      this.token = "";
      this.email = "";
    },
    resetLogin() {
      this.token = "";
      this.email = "";
    },
    async loginRequest(email, passwd) {
      // console.log("jsonData");
      // console.log(jsonData); nojson
      // const url = config.BASE_URL + "/v1/auth/login/passwd";
      const url = config.BASE_URL + "/v1/nojson";
      const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        headers: { "Content-Type": "application/json" },
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
              errorStore.reset();
            } else if (data.code > 400 && data.code < 500) {
              console.log("code > 400 < 500");
              this.resetLogin();
              errorStore.authError(data.message);
            } else {
              console.log("code", data.code);
              errorStore.authError(data.message);
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

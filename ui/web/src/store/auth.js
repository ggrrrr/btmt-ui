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
      console.log("validateRequest validateRequest validateRequest");
      const url = config.BASE_URL + "/v1/auth/validate";
      const requestOptions = {
        method: "POST",
        // headers: this.createHeaders(),
        body: JSON.stringify({}),
      };
      try {
        fetchAPIFunc(url, "payload", requestOptions)
          .then((result) => {
            console.log("result:", result);
          })
          .catch((error) => {
            console.log("error::::", error);
          });
      } catch (error) {
        console.log("catch.error:", error);
      }
    },
    async validateRequest1() {
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
        body: JSON.stringify({
          email: email,
          password: passwd,
        }),
      };
      // try {
      fetchAPIFunc(url, "payload", requestOptions).then(
        (result) => {
          this.email = result.email;
          this.token = result.token;
          console.log("result:", result);
        },
        (error) => {
          // errorStore.networkErr("network error", error.message);
          console.log("error::::", error);
        }
      );
      // } catch (error) {
      // console.log("catch.error:", error);
      // }
    },

    async loginRequest1(email, passwd) {
      const url = config.BASE_URL + "/v1/auth/login/passwd";
      // const url = config.BASE_URL + "/v1/nojson";
      const requestOptions = {
        // mode: "no-cors",
        method: "POST",
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

const fetchAPIFunc = async function (url, payload = "payload", opts = {}) {
  let loginStore = useLoginStore();
  let headers = {
    "Content-Type": "application/json",
  };
  if (loginStore.token) {
    const bearer = "Bearer " + loginStore.token;
    headers["Authorization"] = bearer;
  }

  const options = {
    method: "GET",
    headers: headers,
  };
  if (opts.method) {
    options.method = opts.method;
  }
  if (opts.headers) {
    options.headers.push(opts.method);
  }
  if (opts.body) {
    options.body = opts.body;
  }
  return fetch(url, options).then((response) => {
    if (response.ok) {
      const contentType = response.headers.get("Content-Type") || "";
      if (contentType.includes("application/json")) {
        return response
          .json()
          .then((data) => {
            console.log("json.data:", data);
            if (data[payload]) {
              console.log("json.data.payload:", data[payload]);
              return data[payload];
            }
            console.log("json.data:", data);
            return data;
          })
          .catch((error) => {
            console.log("json.data.error:", error);
            errorStore.systemErr("json error", error);
            return Promise.reject(new Error("Invalid JSON: " + error.message));
          });
      }
      console.log("unkown error:::", response);
      errorStore.systemErr("unkown error", response.statusText);
      return Promise.reject(new Error("unkown error"));
    }

    if (response.status == 401 || response.status == 403) {
      console.log("reset login token", response.status, response.statusText);
      loginStore.resetLogin();
    }

    let message = response.statusText;
    let error = response.statusText;

    return response
      .json()
      .then((data) => {
        console.log("json.error", data);
        console.log("response.status", response.status);
        message = data.message;
        error = data.error ? data.error : "";
        // Auth error
        if (response.status === 401 || response.status === 403) {
          console.log("json.error 401", data);
          errorStore.authError(message);
          return Promise.reject(new Error(`${message} ${error}`));
        }
        if (response.status === 400) {
          console.log("json.error 400", data);

          errorStore.inputErr(message, "");
          return Promise.reject(new Error(`${message} ${error}`));
        }
        console.log("json.error none", data);
        errorStore.networkErr(message, "");
        return Promise.reject(new Error(`${message} ${error}`));
      })
      .catch((err) => {
        console.log("ASDASDASDASDASDASDASD", err);
        // errorStore.networkErr(message, error);
        return Promise.reject(err);
      });
  });
};

export const fetchAPI = fetchAPIFunc;

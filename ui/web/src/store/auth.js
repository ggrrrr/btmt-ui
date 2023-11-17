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
    showLogin: false,
    // error: "",
  }),
  actions: {
    logout() {
      this.token = "";
      this.email = "";
    },
    logedIn(result) {
      this.email = result.email;
      this.token = result.token;
      localStorage.setItem("email", result.email);
      localStorage.setItem("token", result.token);
      this.showLogin = false;
    },
    resetLogin() {
      this.token = "";
      this.email = "";
      this.showLogin = true;
      localStorage.setItem("email", "");
      localStorage.setItem("token", "");
    },
    async validateRequest() {
      const url = config.BASE_URL + "/v1/auth/validate";
      const requestOptions = {
        method: "POST",
        body: JSON.stringify({}),
      };
      const result = await fetchAPIFunc(url, requestOptions);
      if (!result.ok) {
        console.log(result);
      }
    },
    async loginRequest(email, passwd) {
      const url = config.BASE_URL + "/v1/auth/login/passwd";
      const requestOptions = {
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify({
          email: email,
          password: passwd,
        }),
      };
      const { result, ok, error } = await fetchAPIFunc(url, requestOptions);
      if (ok) {
        console.log("result", result);
        this.logedIn(result);
      } else {
        console.log("error:", error);
      }
    },
  },
});

const fetchAPIFunc = function (url, opts = {}) {
  let loginStore = useLoginStore();
  let headers = {
    "Content-Type": "application/json",
  };
  if (loginStore.token) {
    headers["Authorization"] = "Bearer " + loginStore.token;
  }

  const options = {
    method: "GET",
    headers: headers,
  };
  if (opts.method) {
    options.method = opts.method;
  }
  if (opts.headers) {
    options.headers.push(opts.headers);
  }
  if (opts.body) {
    options.body = opts.body;
  }

  let out = {
    result: null,
    ok: false,
    err: null,
  };

  return fetch(url, options)
    .then((response) => {
      console.log("fetch.then OK");
      let message = response.statusText;
      let error = response.statusText;
      return response
        .json()
        .then((data) => {
          if (response.ok) {
            console.log("fetch.then.json OK");
            out.ok = true;
            out.err = null;
            if (data["payload"]) {
              out.result = data["payload"];
              return Promise.resolve(out);
            }
            out.result = data;
            return Promise.resolve(out);
          }
          console.log("fetch.then.json !response.ok");
          console.log("response.status", response.status);
          console.log("error.json.data", data);
          message = data.message;
          error = data.error ? data.error : "";
          // Auth error
          if (response.status === 401 || response.status === 403) {
            console.log("json.error 401/403", data);
            errorStore.authError(message);
            loginStore.resetLogin();

            out.result = null;
            out.ok = false;
            out.err = Error(`${message} ${error}`);
            return Promise.resolve(out);
          }
          if (response.status === 400) {
            console.log("json.error 400", data);
            errorStore.inputErr(message, "");
            out.result = null;
            out.ok = false;
            out.err = Error(`${message} ${error}`);
            return Promise.resolve(out);
          }
          console.log("json.error none", data);
          errorStore.networkErr("Pleasetry again later", response.statusText);
          out.result = null;
          out.ok = false;
          out.err = Error(`${message} ${error}`);
          // return null, false, Error(`${message} ${error}`);
          return Promise.resolve(out);
          // return Promise.reject(new Error(`${message} ${error}`));
        })
        .catch((err) => {
          console.log("catch.json", err);
          errorStore.networkErr(message, error);

          out.result = null;
          out.ok = false;
          out.err = Error(`${message} ${error}`);
          return Promise.resolve(out);
        });
    })
    .catch((error) => {
      console.log("fetch.catch ", error);
      errorStore.networkErr("network error", error);
      out.result = null;
      out.ok = false;
      out.err = error;
      return Promise.resolve(out);
    });
};

export const fetchAPI = fetchAPIFunc;

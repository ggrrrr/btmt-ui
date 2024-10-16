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
    expires_at: localStorage.getItem("expires_at"),
    showLogin: false,
  }),
  actions: {
    logout() {
      this.token = "";
      this.email = "";
      this.expires_at = null;
    },
    loggedIn(result) {
      this.email = result.email;
      this.token = result.token;
      this.expires_at = result.expires_at;
      localStorage.setItem("email", result.email);
      localStorage.setItem("token", result.token);
      localStorage.setItem("expires_at", result.expires_at);
      this.showLogin = false;
    },
    resetLogin() {
      this.token = "";
      this.email = "";
      this.expires_at = null;
      this.showLogin = true;
      localStorage.setItem("email", "");
      localStorage.setItem("token", "");
      localStorage.setItem("expires_at", "");
    },
    async validateRequest() {
      const url = config.BASE_URL + "/auth/token/validate";
      const requestOptions = {
        method: "POST",
        body: JSON.stringify({}),
      };
      const { result, ok, error } = await fetchAPIFunc(url, requestOptions);
      if (!ok) {
        console.log(result);
      } else {
        console.log("error:", error);
      }
    },
    async loginRequest(email, passwd) {
      const url = config.BASE_URL + "/auth/login/passwd";
      const requestOptions = {
        withCredentials: true,
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify({
          email: email,
          password: passwd,
        }),
      };
      const { result, ok, error } = await fetchAPIFunc(url, requestOptions);
      if (ok) {
        let loginResult = {
          email: result.email,
          token: result.token,
          expires_at: parseTimestamp(result.expires_at),
        };
        console.log("result", result);
        let expiresFrom = parseTimestamp(result.expires_at);
        let expires = "expires=" + expiresFrom.toUTCString();
        let cookie = `${result.token};${expires};path=/`;
        document.cookie = "authorization" + "=" + cookie;
        this.loggedIn(loginResult);
      } else {
        console.log("error:", error);
      }
    },
  },
});

const fetchAPIFunc = function (url, opts = {}) {
  let loginStore = useLoginStore();
  let headers = {};
  if (loginStore.token) {
    headers["Authorization"] = "Bearer " + loginStore.token;
  }

  if (opts.formData != undefined) {
    headers["Content-Type"] = "application/json";
  }

  const options = {
    method: "GET",
    headers: headers,
    // credentials: "omit",
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
  if (opts.withCredentials) {
    options.withCredentials = true;
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
      // response.
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

export async function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export function parseTimestamp(fromValue) {
  if (fromValue.seconds) {
    const dateObj = new Date(fromValue.seconds * 1000);
    return dateObj;
  }
  const dateObj = new Date(fromValue);
  return dateObj;
}

export const fetchAPI = fetchAPIFunc;

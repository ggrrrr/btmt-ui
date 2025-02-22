// Utilities
import { defineStore } from "pinia";

import { useConfig } from "@/store/app";
import { useErrorStore } from "@/store/error";
import { apiFetch, parseTimestamp } from "./svc.js";

const config = useConfig;

export class Token {
  constructor(fromJSON) {
    this.value = fromJSON.value;
    this.expires_at = parseTimestamp(fromJSON.expires_at);
  }
}

export function tokenFromJson(fromStr) {
  if (fromStr === undefined) return null;
  if (fromStr === "") return null;

  try {
    let fromJSON = JSON.parse(fromStr);
    if (fromJSON.token === "") return null;
    return new Token(fromJSON);
  } catch (e) {
    return null;
  }
}

let errorStore = useErrorStore();

export const useLoginStore = defineStore({
  id: "auth",
  state: () => ({
    username: localStorage.getItem("username"),
    access_token: tokenFromJson(localStorage.getItem("access_token")),
    refresh_token: tokenFromJson(localStorage.getItem("refresh_token")),
    showLogin: false,
  }),
  actions: {
    logout() {
      this.access_token = null;
      this.refresh_token = null;
      this.username = "";
    },
    loggedIn(result) {
      this.username = result.username;
      this.access_token = result.access_token;
      this.refresh_token = result.refresh_token;
      localStorage.setItem("username", result.username);
      localStorage.setItem("access_token", JSON.stringify(this.access_token));
      localStorage.setItem("refresh_token", JSON.stringify(this.refresh_token));
      let tokenExpires =
        "expires=" + result.access_token.expires_at.toUTCString();
      let tokenCookie = `${result.access_token.value};${tokenExpires};path=/`;
      document.cookie = "authorization" + "=" + tokenCookie;

      this.showLogin = false;
    },
    refreshIn(result) {
      this.access_token = result.access_token;
      localStorage.setItem("access_token", JSON.stringify(this.access_token));
      let tokenExpires =
        "expires=" + result.access_token.expires_at.toUTCString();
      let tokenCookie = `${result.access_token.value};${tokenExpires};path=/`;
      document.cookie = "authorization" + "=" + tokenCookie;

      this.showLogin = false;
    },
    resetLogin() {
      this.username = "";
      this.expires_at = null;
      this.access_token = "";
      this.showLogin = true;
      localStorage.setItem("username", "");
      localStorage.setItem("access_token", "");
      localStorage.setItem("refresh_token", "");
    },
    async validateRequest() {
      const url = config.BASE_URL + "/auth/token/validate";
      const requestOptions = {
        method: "POST",
        body: JSON.stringify({}),
      };
      const { result, ok, error } = await fetchAPI(url, requestOptions);
      if (ok !== true) {
        console.log("validateRequest.error:", error, result);
      }
    },
    async loginRequest(username, passwd) {
      const url = config.BASE_URL + "/v1/auth/login/passwd";
      const requestOptions = {
        // withCredentials: true,
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify({
          username: username,
          password: passwd,
        }),
      };
      const { result, ok, error } = await apiFetch(
        url,
        errorStore,
        this.resetLogin,
        requestOptions
      );
      if (ok) {
        let loginResult = {
          username: result.username,
          access_token: new Token(result.access_token),
          refresh_token: new Token(result.refresh_token),
        };
        this.loggedIn(loginResult);
      } else {
        console.log("loginRequest.error:", error);
      }
    },
  },
});

export const refreshAPI = async function (url, store) {
  const requestOptions = {
    method: "POST",
    token: store.refresh_token.value,
  };

  const { result, ok, error } = await apiFetch(
    url,
    errorStore,
    store.resetLogin,
    requestOptions
  );
  if (ok) {
    let loginResult = {
      username: result.username,
      access_token: new Token(result.access_token),
    };
    store.refreshIn(loginResult);
  } else {
    console.log("refreshAPI.error:", error);
    store.resetLogin();
  }
};

export const fetchAPI = async function (url, opts = {}) {
  let store = useLoginStore();
  if (store.access_token) {
    const now = Date.now();
    const delta = now - store.access_token.expires_at.getTime();
    if (delta > 0) {
      console.log("The token has expired", delta);
      const url = config.BASE_URL + "/v1/auth/token/refresh";

      await refreshAPI(url, store);
    }
    opts.token = store.access_token.value;
  }
  return apiFetch(url, errorStore, store.resetLogin, opts);
};

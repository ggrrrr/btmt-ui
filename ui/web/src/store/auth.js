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
  console.log("tokenFromJson", fromStr);
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
    email: localStorage.getItem("email"),
    access_token: tokenFromJson(localStorage.getItem("access_token")),
    refresh_token: tokenFromJson(localStorage.getItem("refresh_token")),
    showLogin: false,
  }),
  actions: {
    logout() {
      this.access_token = null;
      this.refresh_token = null;
      this.email = "";
    },
    loggedIn(result) {
      this.email = result.email;
      this.access_token = result.access_token;
      this.refresh_token = result.refresh_token;
      localStorage.setItem("email", result.email);
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
      this.email = "";
      this.expires_at = null;
      this.access_token = "";
      this.showLogin = true;
      localStorage.setItem("email", "");
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
      if (!ok) {
        console.log(result);
      } else {
        console.log("error:", error);
      }
    },
    async loginRequest(email, passwd) {
      const url = config.BASE_URL + "/auth/login/passwd";
      const requestOptions = {
        // withCredentials: true,
        // mode: "no-cors",
        method: "POST",
        body: JSON.stringify({
          email: email,
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
          email: result.email,
          access_token: new Token(result.access_token),
          refresh_token: new Token(result.refresh_token),
        };
        console.log("loginResult", loginResult);
        this.loggedIn(loginResult);
      } else {
        console.log("error:", error);
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
      email: result.email,
      access_token: new Token(result.access_token),
    };
    console.log("refreshRequest:", loginResult);
    store.refreshIn(loginResult);
  } else {
    console.log("error:", error);
    store.resetLogin();
  }
};

export const fetchAPI = async function (url, opts = {}) {
  let store = useLoginStore();
  console.log("store", store);
  if (store.access_token) {
    console.log("store.token", store.access_token);
    // let result = store.refreshRequest();
    // console.log("The token is still valid", result);

    const now = Date.now();
    const delta = now - store.access_token.expires_at.getTime();
    console.log("now", now);
    console.log("exp", store.access_token.expires_at.getTime());
    console.log("delta", delta);

    if (delta > 0) {
      console.log("The token has expired", delta);
      const url = config.BASE_URL + "/auth/token/refresh";

      let refreshResponse = await refreshAPI(url, store);
      console.log("refresh", refreshResponse);
    }
    opts.token = store.access_token.value;
  }
  return apiFetch(url, errorStore, store.resetLogin, opts);
};

import { apiFetch } from "./svc.js";

const fakeAuthStore = function () {
  console.log("resetLogin");
};

const fakeErrorStore = {
  authError: (msg, err) => {
    console.log("fakeErrorStore.authError", msg, err);
  },
  inputErr: (msg, err) => {
    console.log("fakeErrorStore.inputErr", msg, err);
  },
  networkErr: (msg, err) => {
    console.log(`fakeErrorStore.networkErr: msg: [${msg}] err: [${err}]`);
  },
};

const result = apiFetch(
  //   "https://google.com/notfound",
  "http://localhost:8010/",
  fakeErrorStore,
  fakeAuthStore,
  {
    token: { token: "asdasd" },
  }
);

console.log("asd", result);

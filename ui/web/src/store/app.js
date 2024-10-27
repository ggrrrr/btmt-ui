// import appConfig from "../app.config.json";

let appConfig = {
  //   BASE_URL: "http://localhost:8010/rest",
};

console.log("location.href", window.location.href);

fetch("/config/config.json")
  .then((resp) => {
    resp
      .json()
      .then((json) => {
        console.log("load json", json);
        appConfig.BASE_URL = json.BASE_URL;
        console.log("appConfig", appConfig);
      })
      .catch((err) => {
        console.log("json error", err);
      });
  })
  .catch((err) => {
    console.log("fetch err", err);
  });

export const useConfig = appConfig;

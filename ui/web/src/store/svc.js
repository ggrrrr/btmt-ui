export const apiFetch = function (url, errorStore, resetLogin, opts = {}) {
  let headers = {};
  if (opts.token) {
    headers["Authorization"] = "Bearer " + opts.token;
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
            console.log("fetch.then.json OK", data);
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
          error = data.error ? data.error : data.message;
          // Auth error
          if (response.status === 401 || response.status === 403) {
            console.log("json.error 401/403", data);
            errorStore.authError(message);
            resetLogin();

            out.result = null;
            out.ok = false;
            out.err = error;
            return Promise.resolve(out);
          }
          if (response.status === 400) {
            console.log("json.error 400", data);
            out.result = null;
            out.ok = false;
            out.err = error;
            errorStore.inputErr(message, error);
            return Promise.resolve(out);
          }
          if (response.status === 404) {
            console.log("json.error 404", data);
            out.result = null;
            out.ok = false;
            out.err = error;
            errorStore.networkErr("Please try again later", error);
            return Promise.resolve(out);
          }
          console.log("json.error unknown", data);
          errorStore.networkErr("Please try again later", response.statusText);
          out.result = null;
          out.ok = false;
          out.err = Error(`${message} ${error}`);
          return Promise.resolve(out);
        })
        .catch((err) => {
          console.log("catch.json", err);
          console.log("catch.json.response", response);
          errorStore.networkErr("response format", err);

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
  if (fromValue === undefined) {
    return null;
  }
  if (fromValue.seconds !== undefined) {
    const dateObj = new Date(fromValue.seconds * 1000);
    return dateObj;
  }

  const dateObj = new Date(fromValue);
  return dateObj;
}

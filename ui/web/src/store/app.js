import appConfig from "../app.config.json";
export const useConfig = appConfig;

// async function parseResponse(res) {
//   // If a body response exists, parse anx extract the possible properties
//   const { data, error, success } = res.status !== 204 ? await res.json() : { success: true };

//   /* Add any custom logic related to how you want to handle the response
//    *
//    * In case success is false,
//    * trigger a new expection to capture later on request call site
//    */
//   if (!success) throw new Error(error.message);
//   // Otherwise, simply resolve the received data
//   return data;
// }

// const apiCallFunc = async function (url, options) {
//   console.log("calling");
//   try {
//     const response = await fetch(url, options);
//     response
//       .json()
//       .then((data) => {
//         console.log("app", data);
//         if (data.code == 200) {
//           console.log("app.data === 200");
//           return data;
//         } else if (data.code == 401 || data.code == 403) {
//           console.log("app.code == 401");
//           //   throw Error(data.message);
//           return Promise.reject(data.messages);
//         } else if (data.code > 400 && data.code < 500) {
//           console.log("app.code > 400 < 500");
//           return Promise.reject(data.messages);
//         } else {
//           console.log("app.code", data.code);
//           return Promise.reject(data.messages);
//         }
//       })
//       .catch((error) => {
//         console.log("app catch", error);
//         // throw Error(error);
//         return Promise.reject(error);
//       });
//   } catch (error) {
//     console.log("app.try.catch");
//     return Promise.reject(error);
//   }
// };

// export const apiCall = apiCallFunc;

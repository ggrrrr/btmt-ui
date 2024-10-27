class Token {
  constructor(fromJSON) {
    console.log("constructor ", fromJSON);

    this.token = fromJSON.token;
    this.expires_at = parseTimestamp(fromJSON.expires_at);
  }
}

const jsonToken = JSON.parse(
  //"ExpiresAt":"2024-10-23T17:28:14.323447Z"}
  //   `{"token":"mytoken","expires_at":"2024-10-23T17:28:14.323447Z"}`

  `{"token":"mytoken","expires_at":{"seconds":1729705992,"nanos":695958000}}`
);

console.log("json ", jsonToken);

// const newToken = Object.assign(new Token(), jsonToken);
const newToken = new Token(jsonToken);

console.log("from json ", newToken);

function parseTimestamp(fromValue) {
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

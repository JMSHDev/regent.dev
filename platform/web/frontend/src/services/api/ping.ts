import { authRequest } from "./auth";

const ping = () => {
  const extraParameters = { params: { id: "PONG" } };
  return authRequest
    .get("/api/ping/", extraParameters)
    .then((response) => Promise.resolve(response.data))
    .catch((error) => Promise.reject(error));
};

export { ping };

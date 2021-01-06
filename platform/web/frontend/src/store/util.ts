export function getTokensFromLocal() {
  const accessToken = localStorage.getItem("accessToken") || "";
  const refreshToken = localStorage.getItem("refreshToken") || "";
  return { accessToken, refreshToken };
}

export function putTokensToLocal(accessToken: string, refreshToken: string) {
  localStorage.setItem("accessToken", accessToken);
  localStorage.setItem("refreshToken", refreshToken);
}

export function delTokensFromLocal() {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
}
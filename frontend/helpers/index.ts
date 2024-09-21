export const urlToWss = (url: string): string => {
  // secure
  if (url.startsWith("https://")) {
    return url.replace(/^https:\/\//, "wss://");
  }

  // insecure
  return url.replace(/^http:\/\//, "ws://");
};
